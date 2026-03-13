package theme

import (
	"embed"
	"fmt"
)

//go:embed css/*.css
var themesFS embed.FS

// Load returns the CSS content for the given theme name ("dark" or "light").
// Note: returns only the theme's scoped variables (under [data-theme="..."]),
// not the base styles. Use LoadAll() for a complete stylesheet.
func Load(name string) (string, error) {
	data, err := themesFS.ReadFile("css/" + name + ".css")
	if err != nil {
		return "", fmt.Errorf("unknown theme %q: %w", name, err)
	}
	return string(data), nil
}

// LoadAll returns the combined CSS for base styles plus both themes.
func LoadAll() (string, error) {
	base, err := themesFS.ReadFile("css/base.css")
	if err != nil {
		return "", fmt.Errorf("load base CSS: %w", err)
	}
	dark, err := themesFS.ReadFile("css/dark.css")
	if err != nil {
		return "", fmt.Errorf("load dark CSS: %w", err)
	}
	light, err := themesFS.ReadFile("css/light.css")
	if err != nil {
		return "", fmt.Errorf("load light CSS: %w", err)
	}
	return string(base) + "\n" + string(dark) + "\n" + string(light), nil
}
