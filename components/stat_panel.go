package components

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
)

type StatPanel struct {
	Title           string
	Description     string
	GridPos         dashboard.GridPos
	Datasource      *DatasourceConfig
	Targets         []*PrometheusQuery
	Transparent     bool
	Mappings        []dashboard.ValueMapping
	Transformations []dashboard.DataTransformerConfig
	Options         *StatPanelOptions
	Thresholds      *dashboard.ThresholdsConfig
	ColorMode       *dashboard.FieldColorModeId
}

type StatPanelOptions struct {
	GraphMode              string
	ColorMode              string
	JustifyMode            string
	TextMode               string
	WideLayout             bool
	ShowPercentChange      bool
	PercentChangeColorMode string
	Orientation            string
	ReduceOptions          *ReduceOptions
}

type ReduceOptions struct {
	Calcs  []string
	Fields string
}

func NewStatPanel(title, description string, gridPos dashboard.GridPos) *StatPanel {
	return &StatPanel{
		Title:           title,
		Description:     description,
		GridPos:         gridPos,
		Transparent:     true,
		Targets:         []*PrometheusQuery{},
		Mappings:        []dashboard.ValueMapping{},
		Transformations: []dashboard.DataTransformerConfig{},
	}
}

func (p *StatPanel) WithDatasource(ds *DatasourceConfig) *StatPanel {
	p.Datasource = ds
	return p
}

func (p *StatPanel) WithTarget(target *PrometheusQuery) *StatPanel {
	p.Targets = append(p.Targets, target)
	return p
}

func (p *StatPanel) WithMappings(mappings []dashboard.ValueMapping) *StatPanel {
	p.Mappings = mappings
	return p
}

func (p *StatPanel) WithTransformations(transformations []dashboard.DataTransformerConfig) *StatPanel {
	p.Transformations = transformations
	return p
}

func (p *StatPanel) WithOptions(options *StatPanelOptions) *StatPanel {
	p.Options = options
	return p
}

// WithThresholds configures custom thresholds for the stat panel
func (p *StatPanel) WithThresholds(mode dashboard.ThresholdsMode, steps []dashboard.Threshold) *StatPanel {
	p.Thresholds = &dashboard.ThresholdsConfig{
		Mode:  mode,
		Steps: steps,
	}
	thresholdsColorMode := dashboard.FieldColorModeId("thresholds")
	p.ColorMode = &thresholdsColorMode

	return p
}

func (p *StatPanel) Build() *stat.PanelBuilder {
	builder := stat.NewPanelBuilder().
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

	for _, target := range p.Targets {
		builder = builder.WithTarget(target.Build())
	}

	if len(p.Mappings) > 0 {
		builder = builder.Mappings(p.Mappings)
	}

	if len(p.Transformations) > 0 {
		builder = builder.Transformations(p.Transformations)
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
	}

	if p.Options != nil {
		reducedDataBuilder := common.NewReduceDataOptionsBuilder()
		reducedDataBuilder.Calcs(p.Options.ReduceOptions.Calcs)
		reducedDataBuilder.Fields(p.Options.ReduceOptions.Fields)
		builder = builder.ReduceOptions(reducedDataBuilder)
	}

	return builder
}
