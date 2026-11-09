package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "resume.md")
	shell := filepath.Join(dir, "template.html")
	if err := os.WriteFile(source, []byte("# A & B\n\n## Experience\n\n- Built **tools** with [Go](https://go.dev).\n\n<script>alert(1)</script>\n\n<!-- private archive -->\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(shell, []byte("<title>{{.Title}}</title><main>{{.Content}}</main>"), 0600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := render(&output, source, shell); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"<title>A &amp; B — Resume</title>", `<h2 id="experience">Experience</h2>`, "<strong>tools</strong>", `<a href="https://go.dev">Go</a>`} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("missing %q in %s", want, output.String())
		}
	}
	for _, unwanted := range []string{"<script>", "private archive"} {
		if strings.Contains(output.String(), unwanted) {
			t.Errorf("raw source leaked: %q", unwanted)
		}
	}
	// Preview requests reread the file, rather than caching the initial source.
	if err := os.WriteFile(source, []byte("# Updated"), 0600); err != nil {
		t.Fatal(err)
	}
	output.Reset()
	if err := render(&output, source, shell); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "<title>Updated — Resume</title>") {
		t.Fatal("source edits were not reflected")
	}
}

func TestRenderFailureDoesNotWritePartialPage(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "resume.md")
	shell := filepath.Join(dir, "template.html")
	if err := os.WriteFile(source, []byte("# Resume"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(shell, []byte("prefix {{.MissingField}}"), 0600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := render(&output, source, shell); err == nil {
		t.Fatal("expected a template error")
	}
	if output.Len() != 0 {
		t.Fatalf("partial output: %q", output.String())
	}
	if err := render(&output, filepath.Join(dir, "missing.md"), shell); err == nil {
		t.Fatal("expected a missing source error")
	}
}
