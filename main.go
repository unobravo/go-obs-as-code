package main

import (
	"log"
	"os"

	"unobravo.com/go-obs-as-code/slo"
)

func main() {
	latencySLO := slo.NewLatencySLO(
		"monthly-agenda-latency-slo-go",
		"Agenda Monthly Latency SLO - 95% requests < 250ms over 28 days",
		"Dashboard to track the monthly latency of the Agenda service: 95% of requests should have latency < 250ms",
		"28d",
		0.95,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_bucket{environment="production", operationName="getDoctorAgenda", job="unobravo-backend", le="250"}`,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_count{environment="production", operationName="getDoctorAgenda", job="unobravo-backend"}`,
	)

	dashboardJSON, err := latencySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating latency dashboard: %v", err)
	}

	outputFile := "output/latency-slo-dashboard.json"
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	if err := os.WriteFile(outputFile, []byte(dashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing latency dashboard file: %v", err)
	}

	availabilitySLO := slo.NewAvailabilitySLO(
		"monthly-agenda-availability-slo-go",
		"Agenda Monthly Availability SLO - 99.9% uptime over 28 days",
		"Dashboard to track the monthly availability of the Agenda service: 99.9% uptime",
		"28d",
		0.999,
		`GraphQL_Errors_total{environment="production",job="unobravo-backend",operationName="getDoctorAgenda",httpStatusCode=~"5.."}`,
		`GraphQL_Requests_total{environment="production",job="unobravo-backend",operationName="getDoctorAgenda"}`,
	)

	availabilityJSON, err := availabilitySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating availability dashboard: %v", err)
	}

	outputFileAvailability := "output/availability-slo-dashboard.json"
	if err := os.WriteFile(outputFileAvailability, []byte(availabilityJSON), 0644); err != nil {
		log.Fatalf("Error writing availability dashboard file: %v", err)
	}
}
