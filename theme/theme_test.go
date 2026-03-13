package theme_test

import (
	"strings"
	"testing"

	"github.com/crosscode-nl/archdiag/theme"
)

func TestLoadDark(t *testing.T) {
	css, err := theme.Load("dark")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if css == "" {
		t.Fatal("dark theme CSS is empty")
	}
	if !strings.Contains(css, `[data-theme="dark"]`) {
		t.Error("dark theme missing [data-theme] selector")
	}
	if !strings.Contains(css, "--bg-primary: #0d1117") {
		t.Error("dark theme missing --bg-primary")
	}
}

func TestLoadLight(t *testing.T) {
	css, err := theme.Load("light")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(css, `[data-theme="light"]`) {
		t.Error("light theme missing [data-theme] selector")
	}
	if !strings.Contains(css, "--bg-primary: #ffffff") {
		t.Error("light theme missing --bg-primary")
	}
}

func TestLoadUnknown(t *testing.T) {
	_, err := theme.Load("neon")
	if err == nil {
		t.Fatal("expected error for unknown theme")
	}
}

func TestLoadAll(t *testing.T) {
	css, err := theme.LoadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if css == "" {
		t.Fatal("combined CSS is empty")
	}
	// Contains base styles
	if !strings.Contains(css, "box-sizing: border-box") {
		t.Error("missing base styles")
	}
	// Contains dark theme variables
	if !strings.Contains(css, `[data-theme="dark"]`) {
		t.Error("missing dark theme selector")
	}
	// Contains light theme variables
	if !strings.Contains(css, `[data-theme="light"]`) {
		t.Error("missing light theme selector")
	}
}
