package slo

import "fmt"

type LatencyQueries struct {
	SuccessMetric string
	TotalMetric   string
	Target        float64
	TimeWindow    string
}

func NewLatencyQueries(successMetric, totalMetric string, target float64, timeWindow string) *LatencyQueries {
	return &LatencyQueries{
		SuccessMetric: successMetric,
		TotalMetric:   totalMetric,
		Target:        target,
		TimeWindow:    timeWindow,
	}
}

func (q *LatencyQueries) SLIQuery() string {
	return fmt.Sprintf(`((sum(rate(%s[5m] offset 2m)) or 0 * sum(rate(%s[5m] offset 2m))) / (sum(rate(%s[5m] offset 2m))))`,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric)
}

func (q *LatencyQueries) SLITimeWindowQuery() string {
	return fmt.Sprintf(`sum(sum_over_time((sum(rate(%s[5m] offset 2m)) or 0 * sum(rate(%s[5m] offset 2m)) < 1e308)[28d:5m])) / sum(sum_over_time((sum(rate(%s[5m] offset 2m)) < 1e308)[28d:5m]))`,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric)
}

func (q *LatencyQueries) FastBurnRateQuery() string {
	return fmt.Sprintf(`(
		(
			(1 - (
				(((sum(rate(%s[5m]) or vector(0))) or 0 * sum(rate(%s[5m]))) / (sum(rate(%s[5m]))))
			)) / (1 - %f) >= 14.4
			and
			(1 - (
				(((sum(rate(%s[1h]) or vector(0))) or 0 * sum(rate(%s[1h]))) / (sum(rate(%s[1h]))))
			)) / (1 - %f) >= 14.4
		)
		or
		(
			(1 - (
				(((sum(rate(%s[30m]) or vector(0))) or 0 * sum(rate(%s[30m]))) / (sum(rate(%s[30m]))))
			)) / (1 - %f) >= 6
			and
			(1 - (
				(((sum(rate(%s[6h]) or vector(0))) or 0 * sum(rate(%s[6h]))) / (sum(rate(%s[6h]))))
			)) / (1 - %f) >= 6
		)
	) or vector(0)`,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target)
}

func (q *LatencyQueries) SlowBurnRateQuery() string {
	return fmt.Sprintf(`(
		(
			(1 - (
				(((sum(rate(%s[2h]) or vector(0))) or 0 * sum(rate(%s[2h]))) / (sum(rate(%s[2h]))))
			)) / (1 - %f) >= 3
			and
			(1 - (
				(((sum(rate(%s[24h]) or vector(0))) or 0 * sum(rate(%s[24h]))) / (sum(rate(%s[24h]))))
			)) / (1 - %f) >= 3
		)
		or
		(
			(1 - (
				(((sum(rate(%s[6h]) or vector(0))) or 0 * sum(rate(%s[6h]))) / (sum(rate(%s[6h]))))
			)) / (1 - %f) >= 1
			and
			(1 - (
				(((sum(rate(%s[72h]) or vector(0))) or 0 * sum(rate(%s[72h]))) / (sum(rate(%s[72h]))))
			)) / (1 - %f) >= 1
		)
	) or vector(0)`,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target)
}

func (q *LatencyQueries) TimeWindowQuery() string {
	return fmt.Sprintf(`label_replace(vector(1), "time_period", "%s", "", "")`, q.TimeWindow)
}

func (q *LatencyQueries) SLOTargetQuery() string {
	return fmt.Sprintf("vector(%f)", q.Target)
}

func (q *LatencyQueries) ErrorBudgetTrendQuery() string {
	return fmt.Sprintf(`((sum(sum_over_time(rate(%s[5m])[%s:4h])) / sum(sum_over_time(rate(%s[5m])[%s:4h]))) - %f) / (1 - %f)`,
		q.SuccessMetric, q.TimeWindow, q.TotalMetric, q.TimeWindow, q.Target, q.Target)
}

func (q *LatencyQueries) RemainingErrorBudgetQuery() string {
	return fmt.Sprintf(`(sum(sum_over_time((sum(rate(%s[5m] offset 2m))< 1e308)[%s:5m])) / sum(sum_over_time((sum(rate(%s[5m] offset 2m))< 1e308
      )[%s:5m])) - %f) / (1 - %f)`,
		q.SuccessMetric, q.TimeWindow, q.TotalMetric, q.TimeWindow, q.Target, q.Target)
}

func (q *LatencyQueries) BurnRateQuery() string {
	return fmt.Sprintf(`avg(1 - avg_over_time(((sum(rate(%s[5m] offset 2m)) / (sum(rate(%s[5m] offset 2m)))) < 1e308)[$__interval:])) / (1 - %f)`,
		q.SuccessMetric, q.TotalMetric, q.Target)
}

func (q *LatencyQueries) InstantBurnRateQuery() string {
	return fmt.Sprintf(`avg(1 - avg_over_time(((sum(rate(%s[5m] offset 2m)) / sum(rate(%s[5m] offset 2m)))< 1e308)[$__interval:])) / (1 - %f)`,
		q.SuccessMetric, q.TotalMetric, q.Target)
}

func (q *LatencyQueries) EventRateQuery() string {
	return fmt.Sprintf(`sum(avg_over_time((sum(rate(%s[5m] offset 2m)))[$__interval:]))`, q.TotalMetric)
}

func (q *LatencyQueries) BurndownFailureEventsQuery() string {
	return fmt.Sprintf(`300 * (sum(sum_over_time(rate(%s[5m])[$__interval:5m] offset 1s)) - sum(sum_over_time(rate(%s[5m])[$__interval:5m] offset 1s)))`,
		q.TotalMetric, q.SuccessMetric)
}

func (q *LatencyQueries) BurndownTotalEventsQuery() string {
	return fmt.Sprintf(`300 * sum(sum_over_time((sum(rate(%s[5m] offset 2m)) < 1e308)[$__range:5m] @ ${__to:date:seconds} offset 1s))`,
		q.TotalMetric)
}
