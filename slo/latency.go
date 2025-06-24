package slo

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"unobravo.com/go-obs-as-code/components"
)

type LatencySLO struct {
	UID                string
	Name               string
	Description        string
	TimeWindow         string
	Target             float64
	SuccessMetricQuery string
	TotalMetricQuery   string
	dashboard          *Dashboard
	queries            *LatencyQueries
}

func NewLatencySLO(uid, name, description, timeWindow string, target float64, successMetricQuery, totalMetricQuery string) *LatencySLO {
	slo := &LatencySLO{
		UID:                uid,
		Name:               name,
		Description:        description,
		TimeWindow:         timeWindow,
		Target:             target,
		SuccessMetricQuery: successMetricQuery,
		TotalMetricQuery:   totalMetricQuery,
		dashboard:          NewDashboard(uid, name, description),
		queries:            NewLatencyQueries(successMetricQuery, totalMetricQuery, target, timeWindow),
	}

	slo.buildRecapRow()
	slo.buildSliRow()
	slo.buildErrorBudgetRow()
	slo.buildBurnRateRow()
	slo.buildEventRateRow()

	return slo
}

func (slo *LatencySLO) BuildJSON() (string, error) {
	return slo.dashboard.ToJSON()
}

// buildRecapRow builds the first row with recap information
func (slo *LatencySLO) buildRecapRow() {
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
		WithOptions(&components.StatPanelOptions{
			GraphMode:              "area",
			ColorMode:              "value",
			JustifyMode:            "auto",
			TextMode:               "auto",
			WideLayout:             true,
			ShowPercentChange:      false,
			PercentChangeColorMode: "standard",
			Orientation:            "auto",
			ReduceOptions: &components.ReduceOptions{
				Calcs:  []string{},
				Fields: "/.*/",
			},
		}).
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
func (slo *LatencySLO) buildSliRow() {
	// SLI timeseries panel
	sliDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("percentunit"),
	}

	sliTarget := components.NewPrometheusQuery("custom_sli", slo.queries.SLIQuery()).WithLegend("SLI")
	sliPanel := components.NewTimeSeriesPanel(
		"SLI",
		"Service level indicator",
		dashboard.GridPos{H: 7, W: 19, X: 0, Y: 4},
	).WithDatasource(sliDS).WithTarget(sliTarget).WithThresholds(dashboard.ThresholdsModeAbsolute, []dashboard.Threshold{
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

	sli28dTarget := components.NewPrometheusQuery("custom_sli_28d", slo.queries.SLITimeWindowQuery()).WithInterval("1m")
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

// buildErrorBudgetRow builds the error budget row
func (slo *LatencySLO) buildErrorBudgetRow() {
	// Error budget trend timeseries
	budgetDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("percentunit"),
	}

	// Primary query for failure events in range
	failureEventsTarget := components.NewPrometheusQuery("Failure in Range", slo.queries.BurndownFailureEventsQuery()).WithLegend("failureEventsInRange")

	// Secondary query for total events
	totalEventsTarget := components.NewPrometheusQuery("Total Events", slo.queries.BurndownTotalEventsQuery()).WithLegend("totalEvents")

	budgetTrendPanel := components.NewTimeSeriesPanel(
		"Error Budget Burndown",
		"The error budget burndown in the selected time\nThe error budget burndown in the selected time",
		dashboard.GridPos{H: 7, W: 19, X: 0, Y: 11},
	).WithDatasource(budgetDS).
		WithTarget(failureEventsTarget).
		WithTarget(totalEventsTarget).
		WithTransformations([]dashboard.DataTransformerConfig{
			{
				Id: "calculateField",
				Options: map[string]interface{}{
					"alias": "cumulativeFailures",
					"cumulative": map[string]interface{}{
						"field":   "failureEventsInRange",
						"reducer": "sum",
					},
					"mode": "cumulativeFunctions",
					"reduce": map[string]interface{}{
						"reducer": "sum",
					},
				},
			},
			{
				Id: "calculateField",
				Options: map[string]interface{}{
					"alias": "totalRemaining",
					"binary": map[string]interface{}{
						"left": map[string]interface{}{
							"matcher": map[string]interface{}{
								"id":      "byName",
								"options": "totalEvents",
							},
						},
						"operator": "-",
						"right": map[string]interface{}{
							"matcher": map[string]interface{}{
								"id":      "byName",
								"options": "cumulativeFailures",
							},
						},
					},
					"mode": "binary",
					"reduce": map[string]interface{}{
						"reducer": "sum",
					},
					"replaceFields": false,
				},
			},
			{
				Id: "calculateField",
				Options: map[string]interface{}{
					"alias": "cumulative sli %",
					"binary": map[string]interface{}{
						"left": map[string]interface{}{
							"matcher": map[string]interface{}{
								"id":      "byName",
								"options": "totalRemaining",
							},
						},
						"operator": "/",
						"right": map[string]interface{}{
							"matcher": map[string]interface{}{
								"id":      "byName",
								"options": "totalEvents",
							},
						},
					},
					"mode": "binary",
					"reduce": map[string]interface{}{
						"reducer": "sum",
					},
					"replaceFields": false,
				},
			},
			{
				Id: "calculateField",
				Options: map[string]interface{}{
					"alias": "sli - objective",
					"binary": map[string]interface{}{
						"left": map[string]interface{}{
							"matcher": map[string]interface{}{
								"id":      "byName",
								"options": "cumulative sli %",
							},
						},
						"operator": "-",
						"right": map[string]interface{}{
							"fixed": fmt.Sprintf("%.6f", slo.Target),
						},
					},
					"mode": "binary",
					"reduce": map[string]interface{}{
						"reducer": "sum",
					},
				},
			},
			{
				Id: "calculateField",
				Options: map[string]interface{}{
					"alias": "error objective",
					"binary": map[string]interface{}{
						"left": map[string]interface{}{
							"fixed": "1",
						},
						"operator": "-",
						"right": map[string]interface{}{
							"fixed": fmt.Sprintf("%.6f", slo.Target),
						},
					},
					"mode": "binary",
					"reduce": map[string]interface{}{
						"reducer": "sum",
					},
				},
			},
			{
				Id: "calculateField",
				Options: map[string]interface{}{
					"alias": "% error budget remaining",
					"binary": map[string]interface{}{
						"left": map[string]interface{}{
							"matcher": map[string]interface{}{
								"id":      "byName",
								"options": "sli - objective",
							},
						},
						"operator": "/",
						"right": map[string]interface{}{
							"matcher": map[string]interface{}{
								"id":      "byName",
								"options": "error objective",
							},
						},
					},
					"mode": "binary",
					"reduce": map[string]interface{}{
						"reducer": "sum",
					},
					"replaceFields": true,
				},
			},
		}).WithThresholds(dashboard.ThresholdsModeAbsolute, []dashboard.Threshold{
		{
			Color: "red",
			Value: float64Ptr(0),
		},
		{
			Color: "yellow",
			Value: float64Ptr(0),
		},
		{
			Color: "green",
			Value: float64Ptr(0.2),
		},
	})
	slo.dashboard.WithPanel(budgetTrendPanel)

	// Remaining error budget
	remainingBudgetDS := &components.DatasourceConfig{
		Type:     "prometheus",
		UID:      "grafanacloud-prom",
		Unit:     stringPtr("percentunit"),
		Min:      float64Ptr(0),
		Max:      float64Ptr(1),
		Decimals: float64Ptr(1),
	}

	remainingBudgetTarget := components.NewPrometheusQuery("custom_remaining_error_budget", slo.queries.RemainingErrorBudgetQuery())
	remainingBudgetPanel := components.NewStatPanel(
		"Remaining Error Budget",
		"The unspent error budget over the last 28d window",
		dashboard.GridPos{H: 7, W: 5, X: 19, Y: 11},
	).WithDatasource(remainingBudgetDS).WithTarget(remainingBudgetTarget).WithThresholds(dashboard.ThresholdsModeAbsolute, []dashboard.Threshold{
		{
			Color: "red",
			Value: float64Ptr(0),
		},
		{
			Color: "yellow",
			Value: float64Ptr(0),
		},
		{
			Color: "green",
			Value: float64Ptr(0.2),
		},
	})
	slo.dashboard.WithPanel(remainingBudgetPanel)
}

// The burn rate row
func (slo *LatencySLO) buildBurnRateRow() {
	// Burn rate timeseries
	burnRateDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("none"),
	}

	burnRateTarget := components.NewPrometheusQuery("custom_burn_rate", slo.queries.BurnRateQuery()).WithLegend("Burn Rate")
	burnRatePanel := components.NewTimeSeriesPanel(
		"Error Budget Burn Rate",
		"The burn rate is the rate that this SLO is spending its error budget over last 5 min [0, 1.0]. A 1x burn rate will consume the entire error budget allotted for that period.",
		dashboard.GridPos{H: 7, W: 19, X: 0, Y: 18},
	).WithDatasource(burnRateDS).WithTarget(burnRateTarget).WithThresholds(dashboard.ThresholdsModeAbsolute, []dashboard.Threshold{
		{
			Color: "green",
			Value: float64Ptr(0),
		},
		{
			Color: "yellow",
			Value: float64Ptr(1),
		},
		{
			Color: "red",
			Value: float64Ptr(3),
		},
	})
	slo.dashboard.WithPanel(burnRatePanel)

	// Current burn rate stat
	currentBurnDS := &components.DatasourceConfig{
		Type:     "prometheus",
		UID:      "grafanacloud-prom",
		Unit:     stringPtr("none"),
		Decimals: float64Ptr(2),
	}

	currentBurnTarget := components.NewPrometheusQuery("current_burn_rate", slo.queries.InstantBurnRateQuery())
	currentBurnPanel := components.NewStatPanel(
		"Current Burn Rate",
		"The burn rate is the rate that this SLO is spending its error budget over last 5 min [0, 1.0]. A 1x burn rate will consume the entire error budget allotted for that period.",
		dashboard.GridPos{H: 7, W: 5, X: 19, Y: 18},
	).WithDatasource(currentBurnDS).WithTarget(currentBurnTarget).WithThresholds(dashboard.ThresholdsModeAbsolute, []dashboard.Threshold{
		{
			Color: "green",
			Value: float64Ptr(0),
		},
		{
			Color: "yellow",
			Value: float64Ptr(1),
		},
		{
			Color: "red",
			Value: float64Ptr(3),
		},
	})
	slo.dashboard.WithPanel(currentBurnPanel)
}

// The event rate row
func (slo *LatencySLO) buildEventRateRow() {
	// Event rate timeseries
	eventRateDS := &components.DatasourceConfig{
		Type: "prometheus",
		UID:  "grafanacloud-prom",
		Unit: stringPtr("reqps"),
	}

	eventRateTarget := components.NewPrometheusQuery("event_rate", slo.queries.EventRateQuery()).WithLegend("Event Rate")
	eventRatePanel := components.NewTimeSeriesPanel(
		"Event Rate",
		"Total Rate (for SLIs that compare rate of successful events to rate of total events, this is the latter)",
		dashboard.GridPos{H: 7, W: 24, X: 0, Y: 25},
	).WithDatasource(eventRateDS).WithTarget(eventRateTarget).WithThresholds(dashboard.ThresholdsModeAbsolute, []dashboard.Threshold{
		{
			Color: "red",
			Value: float64Ptr(80),
		},
		{
			Color: "green",
			Value: float64Ptr(0),
		},
	})
	slo.dashboard.WithPanel(eventRatePanel)
}
