package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cconger/cconger.github.io/resume"
	"github.com/pelletier/go-toml"
)

func render(w io.Writer, dataPath string, templatePath string) error {
	file, err := os.Open(dataPath)
	if err != nil {
		return fmt.Errorf("unable to open datafile %s: %w", dataPath, err)
	}

	var data resume.Resume
	err = toml.NewDecoder(file).Decode(&data)
	if err != nil {
		return fmt.Errorf("unable to process resume data: %w", err)
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("unable to parse template file: %w", err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}
	return nil
}

func liveServe(dataPath string, templatePath string) error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := render(w, dataPath, templatePath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Render error: %v", err)
		}
	}

	s := &http.Server{
		Addr:         ":8080",
		Handler:      http.HandlerFunc(handler),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	return s.ListenAndServe()
}

func main() {
	logger := log.New(os.Stderr, "", 0)

	dataPath := flag.String("data", "", "Required: path to .toml file containing resume")
	templatePath := flag.String("template", "", "Required: path to go template for rendering")
	outputPath := flag.String("output", "", "optional output path otherwise written to stdout")
	serve := flag.Bool("serve", false, "run as a webserver to view output")
	flag.Parse()

	if *dataPath == "" {
		logger.Fatal("data is a required flag")
	}
	if *templatePath == "" {
		logger.Fatal("template is a required flag")
	}

	if *serve {
		err := liveServe(*dataPath, *templatePath)
		if err != nil {
			logger.Fatal("server error: %v", err)
		}
	}

	out := os.Stdout
	if *outputPath != "" {
		// update out to be the file specified here.
		logger.Print("[WARN] File output not yet supported please pipe output")
	}

	err := render(out, *dataPath, *templatePath)
	if err != nil {
		logger.Fatal(err)
	}
}
