package parse_test

import (
	"strings"
	"testing"

	"github.com/crosscode-nl/archdiag/parse"
)

func TestValidateMissingTitle(t *testing.T) {
	yaml := []byte(`
diagram:
  elements:
    - note:
        text: "hello"
`)
	d, err := parse.Parse(yaml)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	errs := parse.Validate(d)
	if len(errs) == 0 {
		t.Fatal("expected validation error for missing title")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "title") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error about title, got: %v", errs)
	}
}

func TestValidateMissingSectionName(t *testing.T) {
	yaml := []byte(`
diagram:
  title: "Test"
  elements:
    - section:
        children:
          - note:
              text: "hello"
`)
	d, err := parse.Parse(yaml)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	errs := parse.Validate(d)
	if len(errs) == 0 {
		t.Fatal("expected validation error for missing section name")
	}
}

func TestValidateMissingFlowName(t *testing.T) {
	yaml := []byte(`
diagram:
  title: "Test"
  elements:
    - flow:
        steps:
          - "A"
`)
	d, err := parse.Parse(yaml)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	errs := parse.Validate(d)
	if len(errs) == 0 {
		t.Fatal("expected validation error for missing flow name")
	}
}

func TestValidateConnectorDirection(t *testing.T) {
	yaml := []byte(`
diagram:
  title: "Test"
  elements:
    - connector:
        direction: diagonal
`)
	d, err := parse.Parse(yaml)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	errs := parse.Validate(d)
	if len(errs) == 0 {
		t.Fatal("expected validation error for invalid direction")
	}
}

func TestValidateValid(t *testing.T) {
	yaml := []byte(`
diagram:
  title: "Valid"
  elements:
    - note:
        text: "hello"
`)
	d, err := parse.Parse(yaml)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	errs := parse.Validate(d)
	if len(errs) != 0 {
		t.Errorf("unexpected validation errors: %v", errs)
	}
}
