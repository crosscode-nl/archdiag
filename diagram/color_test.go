package diagram_test

import (
	"testing"

	"github.com/crosscode-nl/archdiag/diagram"
)

func TestResolveColorHex(t *testing.T) {
	p := diagram.Palette{}
	c, err := diagram.ResolveColor("#4a9eff", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Hex != "#4a9eff" {
		t.Errorf("got Hex=%q, want %q", c.Hex, "#4a9eff")
	}
	if c.Background == "" {
		t.Error("Background should not be empty")
	}
}

func TestResolveColorPalette(t *testing.T) {
	p := diagram.Palette{"aws": "#4a9eff"}
	c, err := diagram.ResolveColor("aws", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Hex != "#4a9eff" {
		t.Errorf("got Hex=%q, want %q", c.Hex, "#4a9eff")
	}
}

func TestResolveColorUnknown(t *testing.T) {
	p := diagram.Palette{"aws": "#4a9eff"}
	_, err := diagram.ResolveColor("unknown", p)
	if err == nil {
		t.Fatal("expected error for unknown color name")
	}
}

func TestResolveColorEmpty(t *testing.T) {
	p := diagram.Palette{}
	c, err := diagram.ResolveColor("", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c != nil {
		t.Errorf("expected nil for empty color, got %+v", c)
	}
}

func TestResolveColorBackground(t *testing.T) {
	p := diagram.Palette{}
	c, err := diagram.ResolveColor("#ff0000", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Background should be rgba with 0.1 opacity
	expected := "rgba(255, 0, 0, 0.1)"
	if c.Background != expected {
		t.Errorf("got Background=%q, want %q", c.Background, expected)
	}
}
