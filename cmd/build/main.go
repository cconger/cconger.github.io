package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func render(w io.Writer, sourcePath, templatePath string) error {
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("read Markdown: %w", err)
	}

	md := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	doc := md.Parser().Parse(text.NewReader(source))
	title := "Resume"
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		if heading, ok := node.(*ast.Heading); ok && heading.Level == 1 {
			title = string(heading.Text(source)) + " — Resume"
			break
		}
	}
	var content bytes.Buffer
	if err := md.Renderer().Render(&content, source, doc); err != nil {
		return fmt.Errorf("render Markdown: %w", err)
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("parse page template: %w", err)
	}
	// Only HTML produced by the Markdown renderer is trusted here. Raw HTML in
	// the source is disabled by Goldmark's default renderer.
	page := struct {
		Title   string
		Content template.HTML
	}{title, template.HTML(content.String())}
	var output bytes.Buffer
	if err := tmpl.Execute(&output, page); err != nil {
		return fmt.Errorf("render page template: %w", err)
	}
	_, err = output.WriteTo(w)
	return err
}

func liveServe(sourcePath, templatePath, addr string) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}
		var page bytes.Buffer
		if err := render(&page, sourcePath, templatePath); err != nil {
			log.Printf("render: %v", err)
			http.Error(w, "Unable to render resume", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = page.WriteTo(w)
	})
	server := &http.Server{
		Addr: addr, Handler: handler,
		ReadTimeout: 3 * time.Second, WriteTimeout: 3 * time.Second,
		IdleTimeout: 30 * time.Second,
	}
	log.Printf("Preview at http://%s (refresh after editing)", addr)
	return server.ListenAndServe()
}

func main() {
	log.SetFlags(0)
	sourcePath := flag.String("data", "resume.md", "path to the Markdown resume")
	templatePath := flag.String("template", "template.html", "path to the page template")
	outputPath := flag.String("output", "", "output file (defaults to stdout)")
	serve := flag.Bool("serve", false, "preview current source on each request")
	addr := flag.String("addr", "127.0.0.1:8080", "preview server address")
	flag.Parse()
	if *serve {
		if err := liveServe(*sourcePath, *templatePath, *addr); err != nil {
			log.Fatal(err)
		}
		return
	}
	var page bytes.Buffer
	if err := render(&page, *sourcePath, *templatePath); err != nil {
		log.Fatal(err)
	}
	if *outputPath != "" {
		if err := os.WriteFile(*outputPath, page.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}
	} else if _, err := page.WriteTo(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
