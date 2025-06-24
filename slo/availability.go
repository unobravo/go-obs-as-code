package slo

import (
	"fmt"
	"time"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"unobravo.com/go-obs-as-code/components"
)

type AvailabilitySLO struct {
	UID                string
	Name               string
	Description        string
	TimeWindow         string
	Target             float64
	SuccessMetricQuery string
	TotalMetricQuery   string
	dashboard          *Dashboard
	queries            *AvailabilityQueries
}

func NewAvailabilitySLO(uid, name, description, timeWindow string, target float64, successMetricQuery, totalMetricQuery string) *AvailabilitySLO {
	slo := &AvailabilitySLO{
		UID:                uid,
		Name:               name,
		Description:        description,
		TimeWindow:         timeWindow,
		Target:             target,
		SuccessMetricQuery: successMetricQuery,
		TotalMetricQuery:   totalMetricQuery,
		dashboard:          NewDashboard(uid, name, description),
		queries:            NewAvailabilityQueries(successMetricQuery, totalMetricQuery, target, timeWindow),
	}

	slo.buildRecapRow()
	slo.buildSliRow()
	slo.buildErrorBudgetRow()
	slo.buildBurnRateRow()
	slo.buildEventRateRow()

	return slo
}

func (slo *AvailabilitySLO) BuildJSON() (string, error) {
	return slo.dashboard.ToJSON()
}

// buildRecapRow builds the first row with recap information
func (slo *AvailabilitySLO) buildRecapRow() {
	// Text panel with title
	textPanel := components.NewTextPanel("", "# "+slo.Name, dashboard.GridPos{H: 4, W: 7, X: 0, Y: 0})
	slo.dashboard.WithPanel(textPanel)

	// Fast burn rate alert panel
	prometheusDS := &components.DatasourceConfig{
		Type:     "prometheus",
		UID:      "grafanacloud-prom",
		Unit:     stringPtr("short"),
		Decimals: float64Ptr(0),
	}

	fastBurnTarget := components.NewPrometheusQuery("fast_burn_alert", slo.queries.FastBurnRateQuery())
	fastBurnPanel := components.NewStatPanel(
		"🚨 Fast Burn Rate Alert",
		"Critical alert when burn rate exceeds fast burn thresholds:\n• 14.4x for 5min AND 1hour\n• 6x for 30min AND 6hour",
		dashboard.GridPos{H: 4, W: 4, X: 7, Y: 0},
	).WithDatasource(prometheusDS).
		WithTarget(fastBurnTarget).
		WithMappings([]dashboard.ValueMapping{
			{
				ValueMap: &dashboard.ValueMap{
					Type: dashboard.MappingTypeValueToText,
					Options: map[string]dashboard.ValueMappingResult{
						"0": {
							Text:  stringPtr("OK"),
							Color: stringPtr("green"),
						},
						"1": {
							Text:  stringPtr("FIRING"),
							Color: stringPtr("red"),
						},
					},
				},
			}})
	slo.dashboard.WithPanel(fastBurnPanel)

	// Slow burn rate alert panel
	slowBurnTarget := components.NewPrometheusQuery("slow_burn_alert", slo.queries.SlowBurnRateQuery())
	slowBurnPanel := components.NewStatPanel(
		"⚠️ Slow Burn Rate Alert",
		"Warning alert for slow burn rate:\n• 3x for 2hours AND 24hours\n• 1x for 6hours AND 72hours",
		dashboard.GridPos{H: 4, W: 4, X: 11, Y: 0},
	).WithDatasource(prometheusDS).
		WithTarget(slowBurnTarget).
		WithMappings([]dashboard.ValueMapping{
			{
				ValueMap: &dashboard.ValueMap{
					Type: dashboard.MappingTypeValueToText,
					Options: map[string]dashboard.ValueMappingResult{
						"0": {
							Text:  stringPtr("OK"),
							Color: stringPtr("green"),
						},
						"1": {
							Text:  stringPtr("CRITICAL"),
							Color: stringPtr("red"),
						},
					},
				},
			}})
	slo.dashboard.WithPanel(slowBurnPanel)

	// Time window panel
	timeWindowDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("percentunit"),
		Min:  float64Ptr(0),
		Max:  float64Ptr(1),
	}

	timeWindowTarget := components.NewPrometheusQuery("time_window", slo.queries.TimeWindowQuery())
	timeWindowPanel := components.NewStatPanel(
		"Time Window",
		"The time window over which the service level objective is being measured over",
		dashboard.GridPos{H: 4, W: 4, X: 15, Y: 0},
	).WithDatasource(timeWindowDS).
		WithTarget(timeWindowTarget).
		WithTransformations([]dashboard.DataTransformerConfig{
			{
				Id: "labelsToFields",
				Options: map[string]interface{}{
					"mode": "rows",
				},
			},
			{
				Id: "organize",
				Options: map[string]interface{}{
					"excludeByName": map[string]interface{}{
						"label": true,
					},
					"indexByName": map[string]interface{}{},
					"renameByName": map[string]interface{}{
						"label": "time_period",
						"value": "Time Window",
					},
				},
			},
		}).WithOptions(&components.StatPanelOptions{
		ReduceOptions: &components.ReduceOptions{
			Calcs:  []string{},
			Fields: "/.*/",
		},
	})

	slo.dashboard.WithPanel(timeWindowPanel)

	// SLO target panel
	sloTargetDS := &components.DatasourceConfig{
		Type:     "prometheus",
		UID:      "grafanacloud-prom",
		Unit:     stringPtr("percentunit"),
		Decimals: float64Ptr(2),
		Min:      float64Ptr(0),
		Max:      float64Ptr(1),
	}

	sloTarget := components.NewPrometheusQuery("A", slo.queries.SLOTargetQuery())
	sloPanel := components.NewStatPanel(
		"SLO",
		"The SLO's Objective value. Always between 0 and 100%",
		dashboard.GridPos{H: 4, W: 5, X: 19, Y: 0},
	).WithDatasource(sloTargetDS).WithTarget(sloTarget)
	slo.dashboard.WithPanel(sloPanel)
}

// SLI row
func (slo *AvailabilitySLO) buildSliRow() {
	// SLI timeseries panel
	sliDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("percentunit"),
	}

	availabilitySLIAvgQuery := fmt.Sprintf(`avg_over_time((
      (((sum(rate(%s[$__rate_interval] offset 2m)) - sum(rate(%s[$__rate_interval] offset 2m) or vector(0))) or 0 * sum(rate(%s[$__rate_interval] offset 2m))) / (sum(rate(%s[$__rate_interval] offset 2m))))
    )[$__interval:])`, slo.TotalMetricQuery, slo.SuccessMetricQuery, slo.TotalMetricQuery, slo.TotalMetricQuery)
	sliTarget1 := components.NewPrometheusQuery("custom_sli_avg", availabilitySLIAvgQuery).WithLegend("AVG")

	futureTimestamp := fmt.Sprintf("%d", time.Now().Unix())
	sliTargetExpr := fmt.Sprintf(`(((sum(rate(%s[$__rate_interval] offset 2m)) - sum(rate(%s[$__rate_interval] offset 2m) or vector(0))) or 0 * sum(rate(%s[$__rate_interval] offset 2m))) / (sum(rate(%s[$__rate_interval] offset 2m)))) AND timestamp(sum(rate(%s[$__rate_interval] offset 2m))) < %s`, slo.TotalMetricQuery, slo.SuccessMetricQuery, slo.TotalMetricQuery, slo.TotalMetricQuery, slo.TotalMetricQuery, futureTimestamp)
	sliTarget2 := components.NewPrometheusQuery("computed_before_creation_time", sliTargetExpr).WithLegend("Before Creation Time")

	sliPanel := components.NewTimeSeriesPanel(
		"SLI",
		"Service level indicator",
		dashboard.GridPos{H: 7, W: 19, X: 0, Y: 4},
	).WithDatasource(sliDS).
		WithTarget(sliTarget1).
		WithTarget(sliTarget2).
		WithThresholds(dashboard.ThresholdsModeAbsolute, []dashboard.Threshold{
			{
				Color: "red",
				Value: float64Ptr(0),
			},
			{
				Color: "green",
				Value: float64Ptr(slo.Target),
			},
		})
	slo.dashboard.WithPanel(sliPanel)

	// SLI 28d stat panel
	sli28dDS := &components.DatasourceConfig{
		Type:     "prometheus",
		UID:      "grafanacloud-prom",
		Unit:     stringPtr("percentunit"),
		Decimals: float64Ptr(1),
		Min:      float64Ptr(0),
		Max:      float64Ptr(1),
	}

	availability28dQuery := fmt.Sprintf(`(
      sum_over_time((sum(rate(%s[5m])) - sum(rate(%s[5m]) or vector(0)))[28d:5m])
      /
      sum_over_time((sum(rate(%s[5m])))[28d:5m])
    )`, slo.TotalMetricQuery, slo.SuccessMetricQuery, slo.TotalMetricQuery)
	sli28dTarget := components.NewPrometheusQuery("custom_sli_28d", availability28dQuery).WithInterval("1m")
	sli28dPanel := components.NewStatPanel(
		"SLI (last 28d)",
		"Service level indicator's value over the last 28d",
		dashboard.GridPos{H: 7, W: 5, X: 19, Y: 4},
	).WithDatasource(sli28dDS).WithTarget(sli28dTarget).WithThresholds(dashboard.ThresholdsModeAbsolute, []dashboard.Threshold{
		{
			Color: "red",
			Value: float64Ptr(0),
		},
		{
			Color: "green",
			Value: float64Ptr(slo.Target),
		},
	})
	slo.dashboard.WithPanel(sli28dPanel)
}

// Error budget row
func (slo *AvailabilitySLO) buildErrorBudgetRow() {
	// Error budget trend timeseries
	budgetDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("percentunit"),
	}

	budgetTrendQuery := fmt.Sprintf(`(
    (
      sum_over_time((sum(rate(%s[5m])) - sum(rate(%s[5m]) or vector(0)))[28d:5m])
      /
      sum_over_time((sum(rate(%s[5m])))[28d:5m])
    ) - %f
  ) / (1 - %f)`, slo.TotalMetricQuery, slo.SuccessMetricQuery, slo.TotalMetricQuery, slo.Target, slo.Target)
	budgetTrendTarget := components.NewPrometheusQuery("custom_error_budget_trend", budgetTrendQuery).WithLegend("Error Budget")
	budgetTrendPanel := components.NewTimeSeriesPanel(
		"Error Budget Trend",
		"If error budget is decreasing over time, it means that your service is spending its error budget faster than it's earning it back.\n\nIf error budget is increasing over time, you're not spending too much of your error budget.",
		dashboard.GridPos{H: 7, W: 19, X: 0, Y: 11},
	).WithDatasource(budgetDS).WithTarget(budgetTrendTarget)
	slo.dashboard.WithPanel(budgetTrendPanel)

	// Remaining error budget stat
	remainingBudgetDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("percentunit"),
		Min:  float64Ptr(1),
		Max:  float64Ptr(1),
	}

	remainingBudgetTarget := components.NewPrometheusQuery("custom_remaining_error_budget", budgetTrendQuery)
	remainingBudgetPanel := components.NewStatPanel(
		"Remaining Error Budget",
		"The unspent error budget over the last 28d window",
		dashboard.GridPos{H: 7, W: 5, X: 19, Y: 11},
	).WithDatasource(remainingBudgetDS).WithTarget(remainingBudgetTarget)
	slo.dashboard.WithPanel(remainingBudgetPanel)
}

// The burn rate row
func (slo *AvailabilitySLO) buildBurnRateRow() {
	// Burn rate timeseries
	burnRateDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("none"),
	}

	// Target 1: AVG
	burnRateAvgQuery := fmt.Sprintf(`(1 - avg_over_time((
      (((sum(rate(%s[5m])) - sum(rate(%s[5m]) or vector(0))) or 0 * sum(rate(%s[5m]))) / (sum(rate(%s[5m]))))
    )[$__interval:])) / (1 - %f)`, slo.TotalMetricQuery, slo.SuccessMetricQuery, slo.TotalMetricQuery, slo.TotalMetricQuery, slo.Target)
	burnRateTarget1 := components.NewPrometheusQuery("custom_burn_rate_avg", burnRateAvgQuery).WithLegend("AVG")

	// Target 2: Instant
	burnRateInstantQuery := fmt.Sprintf(`(1 - (
      (((sum(rate(%s[5m])) - sum(rate(%s[5m]) or vector(0))) or 0 * sum(rate(%s[5m]))) / (sum(rate(%s[5m]))))
    )) / (1 - %f)`, slo.TotalMetricQuery, slo.SuccessMetricQuery, slo.TotalMetricQuery, slo.TotalMetricQuery, slo.Target)
	burnRateTarget2 := components.NewPrometheusQuery("custom_burn_rate_instant", burnRateInstantQuery).WithLegend("Instant")

	burnRatePanel := components.NewTimeSeriesPanel(
		"Error Budget Burn Rate",
		"The burn rate is the rate that this SLO is spending its error budget over last 5 min [0, 1.0]. A 1x burn rate will consume the entire error budget allotted for that period.",
		dashboard.GridPos{H: 7, W: 19, X: 0, Y: 18},
	).WithDatasource(burnRateDS).
		WithTarget(burnRateTarget1).
		WithTarget(burnRateTarget2)
	slo.dashboard.WithPanel(burnRatePanel)

	// Current burn rate
	currentBurnDS := &components.DatasourceConfig{
		Type:     "prometheus",
		UID:      "grafanacloud-prom",
		Unit:     stringPtr("none"),
		Decimals: float64Ptr(2),
	}

	currentBurnTarget := components.NewPrometheusQuery("custom_current_burn_rate", burnRateAvgQuery)
	currentBurnPanel := components.NewStatPanel(
		"Current Burn Rate",
		"The burn rate is the rate that this SLO is spending its error budget over last 5 min [0, 1.0]. A 1x burn rate will consume the entire error budget allotted for that period.",
		dashboard.GridPos{H: 7, W: 5, X: 19, Y: 18},
	).WithDatasource(currentBurnDS).WithTarget(currentBurnTarget)
	slo.dashboard.WithPanel(currentBurnPanel)
}

// The event rate row
func (slo *AvailabilitySLO) buildEventRateRow() {
	// Event rate timeseries
	eventRateDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("reqps"),
	}

	// Target 1: AVG
	eventRateTarget1 := components.NewPrometheusQuery("custom_event_rate", fmt.Sprintf(`sum(rate(%s[$__rate_interval] offset 2m))`, slo.TotalMetricQuery)).WithLegend("AVG")

	// Target 2: Before Creation
	futureTimestamp := fmt.Sprintf("%d", time.Now().Unix())
	eventRateTarget2 := components.NewPrometheusQuery("custom_event_rate_historical", fmt.Sprintf(`sum(rate(%s[$__rate_interval] offset 2m)) AND timestamp(sum(rate(%s[$__rate_interval] offset 2m))) < %s`, slo.TotalMetricQuery, slo.TotalMetricQuery, futureTimestamp)).WithLegend("Before Creation")

	eventRatePanel := components.NewTimeSeriesPanel(
		"Event Rate",
		"Total Rate (for SLIs that compare rate of successful events to rate of total events, this is the latter)",
		dashboard.GridPos{H: 7, W: 24, X: 0, Y: 25},
	).WithDatasource(eventRateDS).
		WithTarget(eventRateTarget1).
		WithTarget(eventRateTarget2)
	slo.dashboard.WithPanel(eventRatePanel)
}

func float64Ptr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}
