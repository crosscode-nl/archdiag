package parse

import (
	"fmt"

	"github.com/crosscode-nl/archdiag/diagram"
	"gopkg.in/yaml.v3"
)

// rawDiagram is the top-level YAML structure.
type rawDiagram struct {
	Diagram struct {
		Title    string            `yaml:"title"`
		Subtitle string            `yaml:"subtitle"`
		Theme    string            `yaml:"theme"`
		Palette  map[string]string `yaml:"palette"`
		Elements []yaml.Node       `yaml:"elements"`
	} `yaml:"diagram"`
}

// Parse parses YAML bytes into a Diagram.
func Parse(data []byte) (*diagram.Diagram, error) {
	var raw rawDiagram
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("yaml parse error: %w", err)
	}

	d := &diagram.Diagram{
		Title:    raw.Diagram.Title,
		Subtitle: raw.Diagram.Subtitle,
		Theme:    raw.Diagram.Theme,
		Palette:  diagram.Palette(raw.Diagram.Palette),
	}

	if d.Theme == "" {
		d.Theme = "dark"
	}
	if d.Palette == nil {
		d.Palette = diagram.Palette{}
	}

	for _, node := range raw.Diagram.Elements {
		comp, err := parseElement(&node, d.Palette)
		if err != nil {
			return nil, err
		}
		d.Elements = append(d.Elements, comp)
	}

	return d, nil
}

// parseElement parses a single YAML mapping node into a Component.
// Each element is a single-key map like {section: {...}} or {flow: {...}}.
func parseElement(node *yaml.Node, palette diagram.Palette) (diagram.Component, error) {
	if node.Kind != yaml.MappingNode || len(node.Content) != 2 {
		return nil, fmt.Errorf("line %d: element must be a single-key mapping", node.Line)
	}

	key := node.Content[0].Value
	value := node.Content[1]

	switch key {
	case "section":
		return parseSection(value, palette)
	case "cards":
		return parseCards(value, palette)
	case "card":
		return parseCard(value, palette)
	case "flow":
		return parseFlow(value, palette)
	case "grid":
		return parseGrid(value, palette)
	case "steps":
		return parseSteps(value, palette)
	case "info":
		return parseInfo(value, palette)
	case "connector":
		return parseConnector(value, palette)
	case "note":
		return parseNote(value)
	default:
		return nil, fmt.Errorf("line %d: unknown element type %q", node.Line, key)
	}
}

func parseSection(node *yaml.Node, palette diagram.Palette) (*diagram.Section, error) {
	var raw struct {
		Name        string      `yaml:"name"`
		Color       string      `yaml:"color"`
		Span        string      `yaml:"span"`
		Border      string      `yaml:"border"`
		Tag         string      `yaml:"tag"`
		Placeholder string      `yaml:"placeholder"`
		Children    []yaml.Node `yaml:"children"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("line %d: section: %w", node.Line, err)
	}

	color, err := diagram.ResolveColor(raw.Color, palette)
	if err != nil {
		return nil, fmt.Errorf("line %d: section %q: %w", node.Line, raw.Name, err)
	}

	s := &diagram.Section{
		Name:        raw.Name,
		Color:       color,
		Span:        raw.Span,
		Border:      raw.Border,
		Tag:         raw.Tag,
		Placeholder: raw.Placeholder,
	}

	for _, childNode := range raw.Children {
		child, err := parseElement(&childNode, palette)
		if err != nil {
			return nil, err
		}
		s.Elements = append(s.Elements, child)
	}

	return s, nil
}

func parseCards(node *yaml.Node, palette diagram.Palette) (*diagram.Cards, error) {
	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("line %d: cards must be a sequence", node.Line)
	}

	cards := &diagram.Cards{}
	for _, itemNode := range node.Content {
		card, err := parseCardItem(itemNode, palette)
		if err != nil {
			return nil, err
		}
		cards.Items = append(cards.Items, *card)
	}
	return cards, nil
}

func parseCard(node *yaml.Node, palette diagram.Palette) (*diagram.Card, error) {
	return parseCardItem(node, palette)
}

func parseCardItem(node *yaml.Node, palette diagram.Palette) (*diagram.Card, error) {
	var raw struct {
		Name     string   `yaml:"name"`
		Color    string   `yaml:"color"`
		Details  []string `yaml:"details"`
		Version  string   `yaml:"version"`
		Subtitle string   `yaml:"subtitle"`
		Footer   string   `yaml:"footer"`
		Span     string   `yaml:"span"`
		Badges   []struct {
			Text  string `yaml:"text"`
			Color string `yaml:"color"`
			Style string `yaml:"style"`
		} `yaml:"badges"`
		Groups []struct {
			Label string   `yaml:"label"`
			Mono  bool     `yaml:"mono"`
			Items []string `yaml:"items"`
		} `yaml:"groups"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("line %d: card: %w", node.Line, err)
	}

	color, err := diagram.ResolveColor(raw.Color, palette)
	if err != nil {
		return nil, fmt.Errorf("line %d: card %q: %w", node.Line, raw.Name, err)
	}

	card := &diagram.Card{
		Name:     raw.Name,
		Color:    color,
		Details:  raw.Details,
		Version:  raw.Version,
		Subtitle: raw.Subtitle,
		Footer:   raw.Footer,
		Span:     raw.Span,
	}

	for _, b := range raw.Badges {
		bc, err := diagram.ResolveColor(b.Color, palette)
		if err != nil {
			return nil, fmt.Errorf("line %d: card %q badge: %w", node.Line, raw.Name, err)
		}
		style := b.Style
		if style == "" {
			style = "filled"
		}
		card.Badges = append(card.Badges, diagram.Badge{
			Text:  b.Text,
			Color: bc,
			Style: style,
		})
	}

	for _, g := range raw.Groups {
		card.Groups = append(card.Groups, diagram.CardGroup{
			Label: g.Label,
			Mono:  g.Mono,
			Items: g.Items,
		})
	}

	return card, nil
}

func parseFlow(node *yaml.Node, palette diagram.Palette) (*diagram.Flow, error) {
	// First try to detect if this is shorthand (steps directly on flow) or full (rows).
	var raw struct {
		Name  string `yaml:"name"`
		Color string `yaml:"color"`
		// Shorthand fields
		Steps     []yaml.Node `yaml:"steps"`
		Connector string      `yaml:"connector"`
		Suffix    string      `yaml:"suffix"`
		// Full rows
		Rows []struct {
			Steps     []yaml.Node `yaml:"steps"`
			Connector string      `yaml:"connector"`
			Suffix    string      `yaml:"suffix"`
		} `yaml:"rows"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("line %d: flow: %w", node.Line, err)
	}

	color, err := diagram.ResolveColor(raw.Color, palette)
	if err != nil {
		return nil, fmt.Errorf("line %d: flow %q: %w", node.Line, raw.Name, err)
	}

	flow := &diagram.Flow{
		Name:  raw.Name,
		Color: color,
	}

	// Shorthand: steps directly on flow
	if len(raw.Steps) > 0 {
		row, err := parseFlowRow(raw.Steps, raw.Connector, raw.Suffix, palette)
		if err != nil {
			return nil, fmt.Errorf("line %d: flow %q: %w", node.Line, raw.Name, err)
		}
		flow.Rows = append(flow.Rows, *row)
	}

	// Full rows
	for _, r := range raw.Rows {
		row, err := parseFlowRow(r.Steps, r.Connector, r.Suffix, palette)
		if err != nil {
			return nil, fmt.Errorf("line %d: flow %q: %w", node.Line, raw.Name, err)
		}
		flow.Rows = append(flow.Rows, *row)
	}

	return flow, nil
}

func parseFlowRow(stepNodes []yaml.Node, connector, suffix string, palette diagram.Palette) (*diagram.FlowRow, error) {
	if connector == "" {
		connector = "arrow"
	}

	row := &diagram.FlowRow{
		Connector: connector,
		Suffix:    suffix,
	}

	for _, sn := range stepNodes {
		step, err := parseFlowStep(&sn, palette)
		if err != nil {
			return nil, err
		}
		row.Steps = append(row.Steps, *step)
	}

	return row, nil
}

func parseFlowStep(node *yaml.Node, palette diagram.Palette) (*diagram.FlowStep, error) {
	// String shorthand: just the text
	if node.Kind == yaml.ScalarNode {
		return &diagram.FlowStep{Text: node.Value}, nil
	}

	var raw struct {
		Text    string   `yaml:"text"`
		Color   string   `yaml:"color"`
		Details []string `yaml:"details"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("line %d: flow step: %w", node.Line, err)
	}

	color, err := diagram.ResolveColor(raw.Color, palette)
	if err != nil {
		return nil, fmt.Errorf("line %d: flow step %q: %w", node.Line, raw.Text, err)
	}

	return &diagram.FlowStep{
		Text:    raw.Text,
		Color:   color,
		Details: raw.Details,
	}, nil
}

func parseGrid(node *yaml.Node, palette diagram.Palette) (*diagram.Grid, error) {
	var raw struct {
		Columns  int         `yaml:"columns"`
		Children []yaml.Node `yaml:"children"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("line %d: grid: %w", node.Line, err)
	}

	grid := &diagram.Grid{Columns: raw.Columns}
	for _, childNode := range raw.Children {
		child, err := parseElement(&childNode, palette)
		if err != nil {
			return nil, err
		}
		grid.Elements = append(grid.Elements, child)
	}

	return grid, nil
}

func parseSteps(node *yaml.Node, palette diagram.Palette) (*diagram.Steps, error) {
	var raw struct {
		Name  string   `yaml:"name"`
		Color string   `yaml:"color"`
		Items []string `yaml:"items"`
		Start int      `yaml:"start"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("line %d: steps: %w", node.Line, err)
	}

	color, err := diagram.ResolveColor(raw.Color, palette)
	if err != nil {
		return nil, fmt.Errorf("line %d: steps %q: %w", node.Line, raw.Name, err)
	}

	start := raw.Start
	if start == 0 {
		start = 1
	}

	return &diagram.Steps{
		Name:  raw.Name,
		Color: color,
		Items: raw.Items,
		Start: start,
	}, nil
}

func parseInfo(node *yaml.Node, palette diagram.Palette) (*diagram.Info, error) {
	var raw struct {
		Name  string `yaml:"name"`
		Color string `yaml:"color"`
		Items []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		} `yaml:"items"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("line %d: info: %w", node.Line, err)
	}

	color, err := diagram.ResolveColor(raw.Color, palette)
	if err != nil {
		return nil, fmt.Errorf("line %d: info %q: %w", node.Line, raw.Name, err)
	}

	info := &diagram.Info{
		Name:  raw.Name,
		Color: color,
	}
	for _, item := range raw.Items {
		info.Items = append(info.Items, diagram.InfoItem{
			Key:   item.Key,
			Value: item.Value,
		})
	}

	return info, nil
}

func parseConnector(node *yaml.Node, palette diagram.Palette) (*diagram.Connector, error) {
	var raw struct {
		Direction string `yaml:"direction"`
		Text      string `yaml:"text"`
		Color     string `yaml:"color"`
		Style     string `yaml:"style"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("line %d: connector: %w", node.Line, err)
	}

	color, err := diagram.ResolveColor(raw.Color, palette)
	if err != nil {
		return nil, fmt.Errorf("line %d: connector: %w", node.Line, err)
	}

	style := raw.Style
	if style == "" {
		style = "arrow"
	}

	return &diagram.Connector{
		Direction: raw.Direction,
		Text:      raw.Text,
		Color:     color,
		Style:     style,
	}, nil
}

func parseNote(node *yaml.Node) (*diagram.Note, error) {
	var raw struct {
		Text  string `yaml:"text"`
		Style string `yaml:"style"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("line %d: note: %w", node.Line, err)
	}

	style := raw.Style
	if style == "" {
		style = "muted"
	}

	return &diagram.Note{
		Text:  raw.Text,
		Style: style,
	}, nil
}
