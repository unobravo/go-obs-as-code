package components

import "github.com/grafana/grafana-foundation-sdk/go/prometheus"

type DatasourceConfig struct {
	Type     string   `json:"type"`
	UID      string   `json:"uid"`
	Unit     *string  `json:"unit,omitempty"`
	Decimals *float64 `json:"decimals,omitempty"`
	Min      *float64 `json:"min,omitempty"`
	Max      *float64 `json:"max,omitempty"`
}

type PrometheusQuery struct {
	RefID        string  `json:"refId"`
	Expr         string  `json:"expr"`
	LegendFormat *string `json:"legendFormat,omitempty"`
	Interval     *string `json:"interval,omitempty"`
	IsRange      bool    `json:"range"`
	IsInstant    bool    `json:"instant"`
}

func NewPrometheusQuery(refID, expr string) *PrometheusQuery {
	return &PrometheusQuery{
		RefID:     refID,
		Expr:      expr,
		IsRange:   true,
		IsInstant: false,
	}
}

func (q *PrometheusQuery) WithLegend(legend string) *PrometheusQuery {
	q.LegendFormat = &legend
	return q
}

func (q *PrometheusQuery) WithInterval(interval string) *PrometheusQuery {
	q.Interval = &interval
	return q
}

func (q *PrometheusQuery) AsRange() *PrometheusQuery {
	q.IsRange = true
	q.IsInstant = false
	return q
}

func (q *PrometheusQuery) AsInstant() *PrometheusQuery {
	q.IsRange = false
	q.IsInstant = true
	return q
}

func (q *PrometheusQuery) Build() *prometheus.DataqueryBuilder {
	builder := prometheus.NewDataqueryBuilder().
		RefId(q.RefID).
		Expr(q.Expr)

	if q.LegendFormat != nil {
		builder = builder.LegendFormat(*q.LegendFormat)
	}

	if q.Interval != nil {
		builder = builder.Interval(*q.Interval)
	}

	if q.IsRange {
		builder = builder.Range()
	}

	if q.IsInstant {
		builder = builder.Instant()
	}

	return builder
}
