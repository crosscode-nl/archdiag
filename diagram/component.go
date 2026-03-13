package diagram

// Component is the interface implemented by all diagram primitives.
type Component interface {
	Type() string
	Children() []Component
}

// Badge is a colored label within a card.
type Badge struct {
	Text  string
	Color *ResolvedColor
	Style string // "filled" (default) or "outlined"
}

// CardGroup is a labeled list of items within a card.
type CardGroup struct {
	Label string
	Mono  bool
	Items []string
}

// FlowStep is a single step within a flow row.
type FlowStep struct {
	Text    string
	Color   *ResolvedColor
	Details []string
}

// FlowRow is a single row of steps within a flow.
type FlowRow struct {
	Steps     []FlowStep
	Connector string // "arrow" (default), "plus", "none"
	Suffix    string
}

// InfoItem is a key-value pair within an info block.
type InfoItem struct {
	Key   string
	Value string
}

// Section is a bordered container with a label. Nestable.
type Section struct {
	Name        string
	Color       *ResolvedColor
	Span        string // "full" (default) or "half"
	Border      string // "solid" (default) or "dashed"
	Tag         string
	Placeholder string
	Elements    []Component
}

func (s *Section) Type() string          { return "section" }
func (s *Section) Children() []Component { return s.Elements }

// Cards is a flex-wrapped row of cards.
type Cards struct {
	Items []Card
}

func (c *Cards) Type() string { return "cards" }
func (c *Cards) Children() []Component {
	out := make([]Component, len(c.Items))
	for i := range c.Items {
		out[i] = &c.Items[i]
	}
	return out
}

// Card is a single card component.
type Card struct {
	Name     string
	Color    *ResolvedColor
	Details  []string
	Badges   []Badge
	Version  string
	Subtitle string
	Footer   string
	Groups   []CardGroup
	Span     string // "full" for grid-column: 1 / -1
}

func (c *Card) Type() string          { return "card" }
func (c *Card) Children() []Component { return nil }

// Flow is a multi-row step sequence with configurable connectors.
type Flow struct {
	Name  string
	Color *ResolvedColor
	Rows  []FlowRow
}

func (f *Flow) Type() string          { return "flow" }
func (f *Flow) Children() []Component { return nil }

// Grid is an explicit multi-column layout.
type Grid struct {
	Columns  int
	Elements []Component
}

func (g *Grid) Type() string          { return "grid" }
func (g *Grid) Children() []Component { return g.Elements }

// Steps is a numbered step list within a named phase.
type Steps struct {
	Name  string
	Color *ResolvedColor
	Items []string
	Start int
}

func (s *Steps) Type() string          { return "steps" }
func (s *Steps) Children() []Component { return nil }

// Info is a key-value grid for metadata.
type Info struct {
	Name  string
	Color *ResolvedColor
	Items []InfoItem
}

func (i *Info) Type() string          { return "info" }
func (i *Info) Children() []Component { return nil }

// Connector is a centered standalone separator element.
type Connector struct {
	Direction string // "down", "right", "up", "left"
	Text      string
	Color     *ResolvedColor
	Style     string // "arrow" (default) or "plain"
}

func (c *Connector) Type() string          { return "connector" }
func (c *Connector) Children() []Component { return nil }

// Note is a footer callout.
type Note struct {
	Text  string
	Style string // "highlight" or "muted" (default)
}

func (n *Note) Type() string          { return "note" }
func (n *Note) Children() []Component { return nil }
