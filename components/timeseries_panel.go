package components

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

type TimeSeriesPanel struct {
	Title           string
	Description     string
	GridPos         dashboard.GridPos
	Datasource      *DatasourceConfig
	Targets         []*PrometheusQuery
	Transparent     bool
	Transformations []dashboard.DataTransformerConfig
	Thresholds      *dashboard.ThresholdsConfig
	ColorMode       *dashboard.FieldColorModeId
	GradientMode    *string
}

func NewTimeSeriesPanel(title, description string, gridPos dashboard.GridPos) *TimeSeriesPanel {
	return &TimeSeriesPanel{
		Title:           title,
		Description:     description,
		GridPos:         gridPos,
		Transparent:     true,
		Targets:         []*PrometheusQuery{},
		Transformations: []dashboard.DataTransformerConfig{},
	}
}

func (p *TimeSeriesPanel) WithDatasource(ds *DatasourceConfig) *TimeSeriesPanel {
	p.Datasource = ds
	return p
}

func (p *TimeSeriesPanel) WithTarget(target *PrometheusQuery) *TimeSeriesPanel {
	p.Targets = append(p.Targets, target)
	return p
}

func (p *TimeSeriesPanel) WithTransformations(transformations []dashboard.DataTransformerConfig) *TimeSeriesPanel {
	p.Transformations = transformations
	return p
}

func (p *TimeSeriesPanel) WithThresholds(mode dashboard.ThresholdsMode, steps []dashboard.Threshold) *TimeSeriesPanel {
	p.Thresholds = &dashboard.ThresholdsConfig{
		Mode:  mode,
		Steps: steps,
	}
	thresholdsColorMode := dashboard.FieldColorModeId("thresholds")
	p.ColorMode = &thresholdsColorMode

	gradientMode := "scheme"
	p.GradientMode = &gradientMode
	return p
}

func (p *TimeSeriesPanel) Build() *timeseries.PanelBuilder {
	builder := timeseries.NewPanelBuilder().
		Title(p.Title).
		Description(p.Description).
		Transparent(p.Transparent).
		GridPos(p.GridPos)

	if p.Datasource != nil {
		builder = builder.Datasource(dashboard.DataSourceRef{
			Type: &p.Datasource.Type,
			Uid:  &p.Datasource.UID,
		})

		if p.Datasource.Unit != nil {
			builder = builder.Unit(*p.Datasource.Unit)
		}
		if p.Datasource.Decimals != nil {
			builder = builder.Decimals(*p.Datasource.Decimals)
		}
		if p.Datasource.Min != nil {
			builder = builder.Min(*p.Datasource.Min)
		}
		if p.Datasource.Max != nil {
			builder = builder.Max(*p.Datasource.Max)
		}
	}

	if p.Thresholds != nil {
		thresholdBuilder := dashboard.NewThresholdsConfigBuilder()
		thresholdBuilder.Mode(p.Thresholds.Mode)
		thresholdBuilder.Steps(p.Thresholds.Steps)
		builder = builder.Thresholds(thresholdBuilder)
	}

	if p.ColorMode != nil {
		colorBuilder := dashboard.NewFieldColorBuilder().Mode(*p.ColorMode)
		builder = builder.ColorScheme(colorBuilder)
		builder.GradientMode(common.GraphGradientMode("scheme"))
	}

	for _, target := range p.Targets {
		builder = builder.WithTarget(target.Build())
	}

	if len(p.Transformations) > 0 {
		builder = builder.Transformations(p.Transformations)
	}

	return builder
}
