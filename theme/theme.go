package theme

import (
	"embed"
	"fmt"
)

//go:embed css/*.css
var themesFS embed.FS

// Load returns the CSS content for the given theme name ("dark" or "light").
func Load(name string) (string, error) {
	data, err := themesFS.ReadFile("css/" + name + ".css")
	if err != nil {
		return "", fmt.Errorf("unknown theme %q: %w", name, err)
	}
	return string(data), nil
}
