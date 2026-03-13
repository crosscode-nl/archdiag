package diagram

import (
	"fmt"
	"strconv"
	"strings"
)

// ResolvedColor holds a hex color and its computed background tint.
type ResolvedColor struct {
	Hex        string // e.g. "#4a9eff"
	Background string // e.g. "rgba(74, 158, 255, 0.1)"
}

// ResolveColor resolves a color string against a palette.
// Returns nil for empty input. Returns error for unknown names.
func ResolveColor(raw string, palette Palette) (*ResolvedColor, error) {
	if raw == "" {
		return nil, nil
	}

	hex := raw
	if !strings.HasPrefix(raw, "#") {
		looked, ok := palette[raw]
		if !ok {
			return nil, fmt.Errorf("unknown color %q", raw)
		}
		hex = looked
	}

	bg, err := hexToRGBA(hex, 0.1)
	if err != nil {
		return nil, fmt.Errorf("invalid hex color %q: %w", hex, err)
	}

	return &ResolvedColor{Hex: hex, Background: bg}, nil
}

// hexToRGBA converts a hex color string to rgba() with the given alpha.
func hexToRGBA(hex string, alpha float64) (string, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "", fmt.Errorf("expected 6 hex digits, got %d", len(hex))
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return "", err
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return "", err
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("rgba(%d, %d, %d, %.1f)", r, g, b, alpha), nil
}
