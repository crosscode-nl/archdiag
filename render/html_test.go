package render_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/crosscode-nl/archdiag/diagram"
	"github.com/crosscode-nl/archdiag/parse"
	"github.com/crosscode-nl/archdiag/render"
)

func TestHTMLRendererMinimal(t *testing.T) {
	d := &diagram.Diagram{
		Title: "Test Diagram",
		Theme: "dark",
		Elements: []diagram.Component{
			&diagram.Note{Text: "Hello, world", Style: "muted"},
		},
	}

	r, err := render.NewHTMLRenderer()
	if err != nil {
		t.Fatalf("NewHTMLRenderer: %v", err)
	}

	var buf bytes.Buffer
	if err := r.Render(d, &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("missing DOCTYPE")
	}
	if !strings.Contains(html, "<title>Test Diagram</title>") {
		t.Error("missing title")
	}
	if !strings.Contains(html, "Hello, world") {
		t.Error("missing note text")
	}
	if !strings.Contains(html, "--bg-primary: #0d1117") {
		t.Error("missing dark theme CSS")
	}
	// Should not contain watch script
	if strings.Contains(html, "EventSource") {
		t.Error("should not contain watch script in normal render")
	}
}

func TestHTMLRendererSection(t *testing.T) {
	d := &diagram.Diagram{
		Title: "Section Test",
		Theme: "dark",
		Elements: []diagram.Component{
			&diagram.Section{
				Name:  "Test Section",
				Color: &diagram.ResolvedColor{Hex: "#4a9eff", Background: "rgba(74, 158, 255, 0.1)"},
				Tag:   "Infra",
				Elements: []diagram.Component{
					&diagram.Note{Text: "nested", Style: "muted"},
				},
			},
		},
	}

	r, err := render.NewHTMLRenderer()
	if err != nil {
		t.Fatalf("NewHTMLRenderer: %v", err)
	}

	var buf bytes.Buffer
	if err := r.Render(d, &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Test Section") {
		t.Error("missing section name")
	}
	if !strings.Contains(html, "#4a9eff") {
		t.Error("missing section color")
	}
	if !strings.Contains(html, "Infra") {
		t.Error("missing section tag")
	}
	if !strings.Contains(html, "nested") {
		t.Error("missing nested note")
	}
}

func TestHTMLRendererLightTheme(t *testing.T) {
	d := &diagram.Diagram{
		Title: "Light Test",
		Theme: "light",
		Elements: []diagram.Component{
			&diagram.Note{Text: "light mode", Style: "muted"},
		},
	}

	r, err := render.NewHTMLRenderer()
	if err != nil {
		t.Fatalf("NewHTMLRenderer: %v", err)
	}

	var buf bytes.Buffer
	if err := r.Render(d, &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "--bg-primary: #ffffff") {
		t.Error("missing light theme CSS")
	}
}

func TestHTMLRendererWatchScript(t *testing.T) {
	d := &diagram.Diagram{
		Title: "Watch Test",
		Theme: "dark",
		Elements: []diagram.Component{
			&diagram.Note{Text: "watch", Style: "muted"},
		},
	}

	r, err := render.NewHTMLRenderer()
	if err != nil {
		t.Fatalf("NewHTMLRenderer: %v", err)
	}

	var buf bytes.Buffer
	if err := r.RenderWithWatch(d, &buf); err != nil {
		t.Fatalf("RenderWithWatch: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "EventSource") {
		t.Error("missing watch script")
	}
}

func TestHTMLRendererAllPrimitives(t *testing.T) {
	data, err := os.ReadFile("../testdata/all-primitives.yaml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	d, err := parse.Parse(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	r, err := render.NewHTMLRenderer()
	if err != nil {
		t.Fatalf("NewHTMLRenderer: %v", err)
	}

	var buf bytes.Buffer
	if err := r.Render(d, &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}

	html := buf.String()

	// Check all primitives rendered
	checks := map[string]string{
		"title":            "All Primitives Test",
		"diagram-subtitle": "Tests every primitive type",
		"section":          "Main Section",
		"nested":           "Nested Section",
		"placeholder":      "No workloads yet",
		"card":             "Card One",
		"badge":            "Filled Badge",
		"version":          "v1.0.0",
		"card-subtitle":    "service-account (namespace)",
		"footer":           "Depends on: Card Two",
		"group":            "s3:GetObject",
		"standalone":       "Standalone Card",
		"flow":             "Request Flow",
		"flow-step":        "Client",
		"flow-detail":      "Maglev hashing",
		"flow-suffix":      "End-to-end encrypted",
		"grid":             "Grid Card 1",
		"steps":            "Phase 1: Infrastructure",
		"steps-item":       "VPC",
		"info":             "Cluster Info",
		"info-key":         "Version",
		"info-value":       "v1.35",
		"connector":        "Internet",
		"note-hl":          "This is a highlighted note",
		"note-muted":       "This is a muted note",
	}

	for name, text := range checks {
		if !strings.Contains(html, text) {
			t.Errorf("%s: expected HTML to contain %q", name, text)
		}
	}
}
