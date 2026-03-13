package diagram

// Palette maps semantic color names to hex values.
type Palette map[string]string

// Diagram is the root of a parsed diagram.
type Diagram struct {
	Title    string
	Subtitle string
	Theme    string // "dark" (default) or "light"
	Palette  Palette
	Elements []Component
}
