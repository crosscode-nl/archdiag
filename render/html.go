package render

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/crosscode-nl/archdiag/diagram"
	"github.com/crosscode-nl/archdiag/theme"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

const watchScript = `new EventSource("/events").onmessage = function() { location.reload(); };`

// HTMLRenderer renders diagrams to self-contained HTML.
type HTMLRenderer struct {
	Tmpl *template.Template
}

// NewHTMLRenderer creates a new HTMLRenderer, parsing all embedded templates.
func NewHTMLRenderer() (*HTMLRenderer, error) {
	r := &HTMLRenderer{}

	// Use the closure pattern: register renderComponent before parsing,
	// closing over the template pointer that will be assigned after ParseFS.
	funcMap := template.FuncMap{
		"renderComponent": func(c diagram.Component) (template.HTML, error) {
			var buf bytes.Buffer
			if err := r.Tmpl.ExecuteTemplate(&buf, c.Type()+".html.tmpl", c); err != nil {
				return "", fmt.Errorf("render %s: %w", c.Type(), err)
			}
			return template.HTML(buf.String()), nil
		},
		"addr": func(c diagram.Card) *diagram.Card {
			return &c
		},
		"add": func(a, b int) int {
			return a + b
		},
		"dirArrow": func(dir string) template.HTML {
			switch dir {
			case "down":
				return "&#x25BC;"
			case "up":
				return "&#x25B2;"
			case "right":
				return "&#x25B6;"
			case "left":
				return "&#x25C0;"
			default:
				return "&#x25BC;"
			}
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	r.Tmpl = tmpl

	return r, nil
}

type renderData struct {
	Diagram      *diagram.Diagram
	ThemeCSS     template.CSS
	DefaultTheme string
	WatchScript  template.JS
}

// Render renders a diagram to HTML without watch script.
func (r *HTMLRenderer) Render(d *diagram.Diagram, w io.Writer) error {
	return r.renderHTML(d, w, false)
}

// RenderWithWatch renders a diagram to HTML with the SSE watch script injected.
func (r *HTMLRenderer) RenderWithWatch(d *diagram.Diagram, w io.Writer) error {
	return r.renderHTML(d, w, true)
}

func (r *HTMLRenderer) renderHTML(d *diagram.Diagram, w io.Writer, watch bool) error {
	css, err := theme.LoadAll()
	if err != nil {
		return fmt.Errorf("load themes: %w", err)
	}

	defaultTheme := d.Theme
	if defaultTheme == "" {
		defaultTheme = "dark"
	}

	data := renderData{
		Diagram:      d,
		ThemeCSS:     template.CSS(css),
		DefaultTheme: defaultTheme,
	}
	if watch {
		data.WatchScript = template.JS(watchScript)
	}

	return r.Tmpl.ExecuteTemplate(w, "diagram.html.tmpl", data)
}
