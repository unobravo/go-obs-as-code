package slo

import "fmt"

type AvailabilityQueries struct {
	SuccessMetric string
	TotalMetric   string
	Target        float64
	TimeWindow    string
}

func NewAvailabilityQueries(successMetric, totalMetric string, target float64, timeWindow string) *AvailabilityQueries {
	return &AvailabilityQueries{
		SuccessMetric: successMetric,
		TotalMetric:   totalMetric,
		Target:        target,
		TimeWindow:    timeWindow,
	}
}

func (q *AvailabilityQueries) SLIQuery() string {
	return fmt.Sprintf(`avg_over_time((
      (((sum(rate(%s[$__rate_interval] offset 2m)) - sum(rate(%s[$__rate_interval] offset 2m) or vector(0))) or 0 * sum(rate(%s[$__rate_interval] offset 2m))) / (sum(rate(%s[$__rate_interval] offset 2m))))
    )[$__interval:])`,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric)
}

func (q *AvailabilityQueries) SLITimeWindowQuery() string {
	return fmt.Sprintf(`(sum_over_time((sum(rate(%s[5m])) - sum(rate(%s[5m]) or vector(0)))[%s:5m]) / sum_over_time((sum(rate(%s[5m])))[%s:5m]))`,
		q.TotalMetric, q.SuccessMetric, q.TimeWindow, q.TotalMetric, q.TimeWindow)
}

func (q *AvailabilityQueries) FastBurnRateQuery() string {
	return fmt.Sprintf(`(
		(
			(1 - (
				(((sum(rate(%s[5m])) - sum(rate(%s[5m]) or vector(0))) or 0 * sum(rate(%s[5m]))) / (sum(rate(%s[5m]))))
			)) / (1 - %f) >= 14.4
			and
			(1 - (
				(((sum(rate(%s[1h])) - sum(rate(%s[1h]) or vector(0))) or 0 * sum(rate(%s[1h]))) / (sum(rate(%s[1h]))))
			)) / (1 - %f) >= 14.4
		)
		or
		(
			(1 - (
				(((sum(rate(%s[30m])) - sum(rate(%s[30m]) or vector(0))) or 0 * sum(rate(%s[30m]))) / (sum(rate(%s[30m]))))
			)) / (1 - %f) >= 6
			and
			(1 - (
				(((sum(rate(%s[6h])) - sum(rate(%s[6h]) or vector(0))) or 0 * sum(rate(%s[6h]))) / (sum(rate(%s[6h]))))
			)) / (1 - %f) >= 6
		)
	) or vector(0)`,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target)
}

func (q *AvailabilityQueries) SlowBurnRateQuery() string {
	return fmt.Sprintf(`(
		(
			(1 - (
				(((sum(rate(%s[2h])) - sum(rate(%s[2h]) or vector(0))) or 0 * sum(rate(%s[2h]))) / (sum(rate(%s[2h]))))
			)) / (1 - %f) >= 3
			and
			(1 - (
				(((sum(rate(%s[24h])) - sum(rate(%s[24h]) or vector(0))) or 0 * sum(rate(%s[24h]))) / (sum(rate(%s[24h]))))
			)) / (1 - %f) >= 3
		)
		or
		(
			(1 - (
				(((sum(rate(%s[6h])) - sum(rate(%s[6h]) or vector(0))) or 0 * sum(rate(%s[6h]))) / (sum(rate(%s[6h]))))
			)) / (1 - %f) >= 1
			and
			(1 - (
				(((sum(rate(%s[72h])) - sum(rate(%s[72h]) or vector(0))) or 0 * sum(rate(%s[72h]))) / (sum(rate(%s[72h]))))
			)) / (1 - %f) >= 1
		)
	) or vector(0)`,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target)
}

func (q *AvailabilityQueries) TimeWindowQuery() string {
	return fmt.Sprintf(`label_replace(vector(1), "time_period", "%s", "", "")`, q.TimeWindow)
}

func (q *AvailabilityQueries) SLOTargetQuery() string {
	return fmt.Sprintf("vector(%f)", q.Target)
}

func (q *AvailabilityQueries) ErrorBudgetTrendQuery() string {
	return fmt.Sprintf(`((sum_over_time((sum(rate(%s[5m])) - sum(rate(%s[5m]) or vector(0)))[%s:5m]) / sum_over_time((sum(rate(%s[5m])))[%s:5m])) - %f) / (1 - %f)`,
		q.TotalMetric, q.SuccessMetric, q.TimeWindow, q.TotalMetric, q.TimeWindow, q.Target, q.Target)
}

func (q *AvailabilityQueries) RemainingErrorBudgetQuery() string {
	return q.ErrorBudgetTrendQuery()
}

func (q *AvailabilityQueries) BurnRateQuery() string {
	return fmt.Sprintf(`(1 - avg_over_time((
      (((sum(rate(%s[5m])) - sum(rate(%s[5m]) or vector(0))) or 0 * sum(rate(%s[5m]))) / (sum(rate(%s[5m]))))
    )[$__interval:])) / (1 - %f)`,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target)
}

func (q *AvailabilityQueries) InstantBurnRateQuery() string {
	return fmt.Sprintf(`(1 - avg_over_time((
      (((sum(rate(%s[5m])) - sum(rate(%s[5m]) or vector(0))) or 0 * sum(rate(%s[5m]))) / (sum(rate(%s[5m]))))
    )[$__interval:])) / (1 - %f)`,
		q.TotalMetric, q.SuccessMetric, q.TotalMetric, q.TotalMetric, q.Target)
}

func (q *AvailabilityQueries) EventRateQuery() string {
	return fmt.Sprintf(`sum(rate(%s[$__rate_interval] offset 2m))`, q.TotalMetric)
}
