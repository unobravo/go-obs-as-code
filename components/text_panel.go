package components

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/text"
)

type TextPanel struct {
	Title       string
	Content     string
	GridPos     dashboard.GridPos
	Transparent bool
}

func NewTextPanel(title, content string, gridPos dashboard.GridPos) *TextPanel {
	return &TextPanel{
		Title:       title,
		Content:     content,
		GridPos:     gridPos,
		Transparent: true,
	}
}

func (p *TextPanel) Build() *text.PanelBuilder {
	return text.NewPanelBuilder().
		Title(p.Title).
		Content(p.Content).
		Transparent(p.Transparent).
		GridPos(p.GridPos)
}
