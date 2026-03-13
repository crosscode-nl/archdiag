package render

import (
	"io"

	"github.com/crosscode-nl/archdiag/diagram"
)

// Renderer renders a Diagram to an output stream.
type Renderer interface {
	Render(d *diagram.Diagram, w io.Writer) error
}
