package parse_test

import (
	"os"
	"testing"

	"github.com/crosscode-nl/archdiag/parse"
)

func TestParseMinimal(t *testing.T) {
	data, err := os.ReadFile("../testdata/minimal.yaml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	d, err := parse.Parse(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if d.Title != "Minimal Test" {
		t.Errorf("Title=%q, want %q", d.Title, "Minimal Test")
	}
	if d.Subtitle != "A minimal diagram" {
		t.Errorf("Subtitle=%q, want %q", d.Subtitle, "A minimal diagram")
	}
	if d.Theme != "dark" {
		t.Errorf("Theme=%q, want %q (default)", d.Theme, "dark")
	}
	if len(d.Elements) != 1 {
		t.Fatalf("Elements count=%d, want 1", len(d.Elements))
	}
	if d.Elements[0].Type() != "note" {
		t.Errorf("Element[0].Type()=%q, want %q", d.Elements[0].Type(), "note")
	}
}

func TestParsePalette(t *testing.T) {
	yaml := []byte(`
diagram:
  title: "Palette Test"
  palette:
    aws: "#4a9eff"
    azure: "#a855f7"
  elements:
    - section:
        name: "AWS"
        color: aws
        children:
          - note:
              text: "inside"
`)
	d, err := parse.Parse(yaml)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(d.Palette) != 2 {
		t.Errorf("Palette count=%d, want 2", len(d.Palette))
	}
	section := d.Elements[0]
	if section.Type() != "section" {
		t.Fatalf("Element[0].Type()=%q, want section", section.Type())
	}
}

func TestParseFlowShorthand(t *testing.T) {
	yaml := []byte(`
diagram:
  title: "Flow Shorthand"
  elements:
    - flow:
        name: "Test Flow"
        steps:
          - "Step A"
          - { text: "Step B", color: "#ff0000" }
`)
	d, err := parse.Parse(yaml)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if d.Elements[0].Type() != "flow" {
		t.Fatalf("expected flow, got %s", d.Elements[0].Type())
	}
}

func TestParseAllPrimitives(t *testing.T) {
	yaml := []byte(`
diagram:
  title: "All Primitives"
  palette:
    blue: "#4a9eff"
  elements:
    - section:
        name: "Test Section"
        color: blue
        border: dashed
        tag: "Infrastructure"
        children:
          - cards:
              - name: "Card A"
                details: ["detail 1"]
                badges:
                  - { text: "badge", color: blue, style: outlined }
                version: "v1.0"
                subtitle: "sub"
                footer: "Depends on: X"
                groups:
                  - label: "Perms"
                    mono: true
                    items: ["perm1", "perm2"]
    - card:
        name: "Standalone Card"
        span: full
    - flow:
        name: "Test Flow"
        color: blue
        rows:
          - steps:
              - { text: "A", color: blue, details: ["sub-detail"] }
              - "B"
            connector: plus
            suffix: "suffix text"
          - steps:
              - "C"
              - "D"
            connector: arrow
    - grid:
        columns: 2
        children:
          - card:
              name: "Grid Card 1"
          - card:
              name: "Grid Card 2"
    - steps:
        name: "Phase 1"
        color: blue
        start: 5
        items: ["Step A", "Step B"]
    - info:
        name: "Metadata"
        color: blue
        items:
          - { key: "Key1", value: "Val1" }
    - connector:
        direction: down
        text: "Internet"
        color: blue
        style: arrow
    - note:
        text: "A note"
        style: highlight
`)
	d, err := parse.Parse(yaml)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	expectedTypes := []string{"section", "card", "flow", "grid", "steps", "info", "connector", "note"}
	if len(d.Elements) != len(expectedTypes) {
		t.Fatalf("Elements count=%d, want %d", len(d.Elements), len(expectedTypes))
	}
	for i, et := range expectedTypes {
		if d.Elements[i].Type() != et {
			t.Errorf("Element[%d].Type()=%q, want %q", i, d.Elements[i].Type(), et)
		}
	}
}
