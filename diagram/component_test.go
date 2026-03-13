package diagram_test

import (
	"testing"

	"github.com/crosscode-nl/archdiag/diagram"
)

// TestComponentInterface verifies all concrete types implement Component.
func TestComponentInterface(t *testing.T) {
	components := []diagram.Component{
		&diagram.Section{Name: "test"},
		&diagram.Cards{Items: nil},
		&diagram.Card{Name: "test"},
		&diagram.Flow{Name: "test"},
		&diagram.Grid{Columns: 2},
		&diagram.Steps{Name: "test"},
		&diagram.Info{Name: "test"},
		&diagram.Connector{Direction: "down"},
		&diagram.Note{Text: "test"},
	}

	expectedTypes := []string{
		"section", "cards", "card", "flow", "grid", "steps", "info", "connector", "note",
	}

	for i, c := range components {
		if c.Type() != expectedTypes[i] {
			t.Errorf("component %d: got Type()=%q, want %q", i, c.Type(), expectedTypes[i])
		}
	}
}

func TestLeafChildrenEmpty(t *testing.T) {
	leaves := []diagram.Component{
		&diagram.Card{Name: "test"},
		&diagram.Steps{Name: "test"},
		&diagram.Info{Name: "test"},
		&diagram.Connector{Direction: "down"},
		&diagram.Note{Text: "test"},
	}
	for _, c := range leaves {
		if ch := c.Children(); len(ch) != 0 {
			t.Errorf("%s.Children() returned %d items, want 0", c.Type(), len(ch))
		}
	}
}

func TestSectionChildren(t *testing.T) {
	child := &diagram.Card{Name: "child"}
	s := &diagram.Section{
		Name:     "parent",
		Elements: []diagram.Component{child},
	}
	if len(s.Children()) != 1 {
		t.Fatalf("Section.Children() returned %d items, want 1", len(s.Children()))
	}
	if s.Children()[0].Type() != "card" {
		t.Errorf("got child type %q, want %q", s.Children()[0].Type(), "card")
	}
}
