package components

import "github.com/charmbracelet/glamour"

// MarkdownRenderer defines the interface for rendering markdown content
type MarkdownRenderer interface {
	Render(markdown string, width int) string
}

// DefaultRenderer implements MarkdownRenderer with standard styling
type DefaultRenderer struct {
	renderer *glamour.TermRenderer
	width    int
}

// NewDefaultRenderer creates a new DefaultRenderer instance
func NewDefaultRenderer() *DefaultRenderer {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(80),
	)
	return &DefaultRenderer{
		renderer: renderer,
		width:    80,
	}
}

// Render renders markdown content with standard styling
func (r *DefaultRenderer) Render(markdown string, width int) string {
	if r.renderer == nil || width != r.width {
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithEmoji(),
			glamour.WithWordWrap(width),
		)
		if err != nil {
			return "Error creating renderer: " + err.Error()
		}
		r.renderer = renderer
		r.width = width
	}

	rendered, err := r.renderer.Render(markdown)
	if err != nil {
		return "Error rendering markdown: " + err.Error()
	}

	return rendered
}

// TestRenderer implements MarkdownRenderer with test-friendly styling
type TestRenderer struct {
	renderer *glamour.TermRenderer
	width    int
}

// NewTestRenderer creates a new TestRenderer instance
func NewTestRenderer() *TestRenderer {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("notty"),
		glamour.WithWordWrap(80),
	)
	return &TestRenderer{
		renderer: renderer,
		width:    80,
	}
}

// Render renders markdown content without terminal styling
func (r *TestRenderer) Render(markdown string, width int) string {
	if r.renderer == nil || width != r.width {
		renderer, err := glamour.NewTermRenderer(
			glamour.WithStandardStyle("notty"),
			glamour.WithWordWrap(width),
		)
		if err != nil {
			return "Error creating renderer: " + err.Error()
		}
		r.renderer = renderer
		r.width = width
	}

	rendered, err := r.renderer.Render(markdown)
	if err != nil {
		return "Error rendering markdown: " + err.Error()
	}

	return rendered
}
