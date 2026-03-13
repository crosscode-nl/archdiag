package parse

import (
	"fmt"

	"github.com/crosscode-nl/archdiag/diagram"
)

// Validate checks a parsed Diagram for semantic errors.
// Returns a slice of errors (empty if valid).
func Validate(d *diagram.Diagram) []error {
	var errs []error

	if d.Title == "" {
		errs = append(errs, fmt.Errorf("diagram: title is required"))
	}

	for i, elem := range d.Elements {
		errs = append(errs, validateComponent(elem, fmt.Sprintf("elements[%d]", i))...)
	}

	return errs
}

func validateComponent(c diagram.Component, path string) []error {
	var errs []error

	switch v := c.(type) {
	case *diagram.Section:
		if v.Name == "" {
			errs = append(errs, fmt.Errorf("%s: section name is required", path))
		}
		if v.Border != "" && v.Border != "solid" && v.Border != "dashed" {
			errs = append(errs, fmt.Errorf("%s: section border must be 'solid' or 'dashed', got %q", path, v.Border))
		}
		if v.Span != "" && v.Span != "full" && v.Span != "half" {
			errs = append(errs, fmt.Errorf("%s: section span must be 'full' or 'half', got %q", path, v.Span))
		}
		for i, child := range v.Elements {
			errs = append(errs, validateComponent(child, fmt.Sprintf("%s.section(%s).children[%d]", path, v.Name, i))...)
		}

	case *diagram.Cards:
		for i, card := range v.Items {
			if card.Name == "" {
				errs = append(errs, fmt.Errorf("%s.cards[%d]: card name is required", path, i))
			}
		}

	case *diagram.Card:
		if v.Name == "" {
			errs = append(errs, fmt.Errorf("%s: card name is required", path))
		}

	case *diagram.Flow:
		if v.Name == "" {
			errs = append(errs, fmt.Errorf("%s: flow name is required", path))
		}
		for i, row := range v.Rows {
			if row.Connector != "arrow" && row.Connector != "plus" && row.Connector != "none" {
				errs = append(errs, fmt.Errorf("%s.flow(%s).rows[%d]: connector must be 'arrow', 'plus', or 'none', got %q", path, v.Name, i, row.Connector))
			}
			for j, step := range row.Steps {
				if step.Text == "" {
					errs = append(errs, fmt.Errorf("%s.flow(%s).rows[%d].steps[%d]: step text is required", path, v.Name, i, j))
				}
			}
		}

	case *diagram.Grid:
		if v.Columns < 1 {
			errs = append(errs, fmt.Errorf("%s: grid columns must be >= 1", path))
		}
		for i, child := range v.Elements {
			errs = append(errs, validateComponent(child, fmt.Sprintf("%s.grid.children[%d]", path, i))...)
		}

	case *diagram.Steps:
		if v.Name == "" {
			errs = append(errs, fmt.Errorf("%s: steps name is required", path))
		}
		if len(v.Items) == 0 {
			errs = append(errs, fmt.Errorf("%s: steps must have at least one item", path))
		}

	case *diagram.Info:
		if v.Name == "" {
			errs = append(errs, fmt.Errorf("%s: info name is required", path))
		}

	case *diagram.Connector:
		validDirs := map[string]bool{"down": true, "right": true, "up": true, "left": true}
		if !validDirs[v.Direction] {
			errs = append(errs, fmt.Errorf("%s: connector direction must be down/right/up/left, got %q", path, v.Direction))
		}
		if v.Style != "arrow" && v.Style != "plain" {
			errs = append(errs, fmt.Errorf("%s: connector style must be 'arrow' or 'plain', got %q", path, v.Style))
		}

	case *diagram.Note:
		if v.Text == "" {
			errs = append(errs, fmt.Errorf("%s: note text is required", path))
		}
		if v.Style != "highlight" && v.Style != "muted" {
			errs = append(errs, fmt.Errorf("%s: note style must be 'highlight' or 'muted', got %q", path, v.Style))
		}
	}

	return errs
}
