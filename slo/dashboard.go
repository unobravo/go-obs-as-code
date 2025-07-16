package slo

import (
	"encoding/json"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"unobravo.com/go-obs-as-code/components"
)

// Grafana dashboard with panels
type Dashboard struct {
	UID         string
	Title       string
	Description string
	builder     *dashboard.DashboardBuilder
}

func NewDashboard(uid, title, description string) *Dashboard {
	builder := dashboard.NewDashboardBuilder(title).
		Uid(uid).
		Description(description).
		Editable().
		Time("now-6h", "now").
		Tags([]string{"slo"})

	return &Dashboard{
		UID:         uid,
		Title:       title,
		Description: description,
		builder:     builder,
	}
}

func (d *Dashboard) WithPanel(panel interface{}) *Dashboard {
	switch p := panel.(type) {
	case *components.StatPanel:
		d.builder = d.builder.WithPanel(p.Build())
	case *components.TimeSeriesPanel:
		d.builder = d.builder.WithPanel(p.Build())
	case *components.TextPanel:
		d.builder = d.builder.WithPanel(p.Build())
	}
	return d
}

func (d *Dashboard) Build() (*dashboard.DashboardBuilder, error) {
	return d.builder, nil
}

func (d *Dashboard) ToJSON() (string, error) {
	dashboard, err := d.builder.Build()
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.MarshalIndent(dashboard, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
