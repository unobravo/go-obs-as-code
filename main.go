package main

import (
	"log"
	"os"

	"unobravo.com/go-obs-as-code/slo"
)

func main() {
	latencySLO := slo.NewLatencySLO(
		"monthly-agenda-latency-slo",
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

	freeAppointmentCreationLatencySLO := slo.NewLatencySLO(
		"free-appointment-creation-latency-slo",
		"Free Appointment Creation Latency SLO - 95% requests < 350ms over 28 days",
		"Dashboard to track the monthly latency of the Free Appointment Creation service: 95% of requests should have latency < 350ms",
		"28d",
		0.95,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_bucket{environment="production", operationName="createSessionByPatient", job="unobravo-backend", le="350"}`,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_count{environment="production", operationName="createSessionByPatient", job="unobravo-backend"}`,
	)

	freeAppointmentCreationDashboardJSON, err := freeAppointmentCreationLatencySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating latency dashboard: %v", err)
	}

	outputFileFreeAppointmentCreation := "output/free-appointment-creation-latency-slo-dashboard.json"
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	if err := os.WriteFile(outputFileFreeAppointmentCreation, []byte(freeAppointmentCreationDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing latency dashboard file: %v", err)
	}

	freeAppointmentCreationAvailabilitySLO := slo.NewAvailabilitySLO(
		"free-session-creation-availability-slo",
		"Free Appointment Creation Availability SLO - 99.9% uptime over 28 days",
		"Dashboard to track the monthly availability of the Free Appointment Creation service: 99.9% uptime",
		"28d",
		0.999,
		`GraphQL_Errors_total{environment="production",job="unobravo-backend",operationName="createSessionByPatient",httpStatusCode=~"5.."}`,
		`GraphQL_Requests_total{environment="production",job="unobravo-backend",operationName="createSessionByPatient"}`,
	)

	freeAppointmentCreationAvailabilityDashboardJSON, err := freeAppointmentCreationAvailabilitySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating availability dashboard: %v", err)
	}

	outputFileFreeAppointmentCreationAvailability := "output/free-appointment-creation-availability-slo-dashboard.json"
	if err := os.WriteFile(outputFileFreeAppointmentCreationAvailability, []byte(freeAppointmentCreationAvailabilityDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing availability dashboard file: %v", err)
	}

	freeAppointmentUpdateAvailabilitySLO := slo.NewAvailabilitySLO(
		"free-session-update-availability-slo",
		"Free Appointment Update Availability SLO - 99.9% uptime over 28 days",
		"Dashboard to track the monthly availability of the Free Appointment Update service: 99.9% uptime",
		"28d",
		0.999,
		`GraphQL_Errors_total{environment="production",job="unobravo-backend",operationName="updateSessionByPatient",httpStatusCode=~"5.."}`,
		`GraphQL_Requests_total{environment="production",job="unobravo-backend",operationName="updateSessionByPatient"}`,
	)

	freeAppointmentUpdateAvailabilityDashboardJSON, err := freeAppointmentUpdateAvailabilitySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating availability dashboard: %v", err)
	}

	outputFileFreeAppointmentUpdateAvailability := "output/free-appointment-update-availability-slo-dashboard.json"
	if err := os.WriteFile(outputFileFreeAppointmentUpdateAvailability, []byte(freeAppointmentUpdateAvailabilityDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing availability dashboard file: %v", err)
	}

	freeAppointmentUpdateLatencySLO := slo.NewLatencySLO(
		"free-session-update-latency-slo",
		"Free Appointment Update Latency SLO - 95% requests < 350ms over 28 days",
		"Dashboard to track the monthly latency of the Free Appointment Update service: 95% of requests should have latency < 350ms",
		"28d",
		0.95,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_bucket{environment="production", operationName="updateSessionByPatient", job="unobravo-backend", le="350"}`,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_count{environment="production", operationName="updateSessionByPatient", job="unobravo-backend"}`,
	)

	freeAppointmentUpdateLatencyDashboardJSON, err := freeAppointmentUpdateLatencySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating latency dashboard: %v", err)
	}

	outputFileFreeAppointmentUpdateLatency := "output/free-appointment-update-latency-slo-dashboard.json"
	if err := os.WriteFile(outputFileFreeAppointmentUpdateLatency, []byte(freeAppointmentUpdateLatencyDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing latency dashboard file: %v", err)
	}

	freeSessionDeleteAvailabilitySLO := slo.NewAvailabilitySLO(
		"free-session-delete-availability-slo",
		"Free Session Delete Availability SLO - 99.9% uptime over 28 days",
		"Dashboard to track the monthly availability of the Free Session Delete service: 99.9% uptime",
		"28d",
		0.999,
		`GraphQL_Errors_total{environment="production",job="unobravo-backend",operationName="cancelSessionByPatient",httpStatusCode=~"5.."}`,
		`GraphQL_Requests_total{environment="production",job="unobravo-backend",operationName="cancelSessionByPatient"}`,
	)

	freeSessionDeleteAvailabilityDashboardJSON, err := freeSessionDeleteAvailabilitySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating availability dashboard: %v", err)
	}

	outputFileFreeSessionDeleteAvailability := "output/free-session-delete-availability-slo-dashboard.json"
	if err := os.WriteFile(outputFileFreeSessionDeleteAvailability, []byte(freeSessionDeleteAvailabilityDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing availability dashboard file: %v", err)
	}

	freeSessionDeleteLatencySLO := slo.NewLatencySLO(
		"free-session-delete-latency-slo",
		"Free Session Delete Latency SLO - 95% requests < 350ms over 28 days",
		"Dashboard to track the monthly latency of the Free Session Delete service: 95% of requests should have latency < 350ms",
		"28d",
		0.95,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_bucket{environment="production", operationName="cancelSessionByPatient", job="unobravo-backend", le="350"}`,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_count{environment="production", operationName="cancelSessionByPatient", job="unobravo-backend"}`,
	)

	freeSessionDeleteLatencyDashboardJSON, err := freeSessionDeleteLatencySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating latency dashboard: %v", err)
	}

	outputFileFreeSessionDeleteLatency := "output/free-session-delete-latency-slo-dashboard.json"
	if err := os.WriteFile(outputFileFreeSessionDeleteLatency, []byte(freeSessionDeleteLatencyDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing latency dashboard file: %v", err)
	}

	getConversationsAvailabilitySLO := slo.NewAvailabilitySLO(
		"get-conversations-availability-slo",
		"Get Conversations Availability SLO - 99.9% uptime over 28 days",
		"Dashboard to track the monthly availability of the Get Conversations service: 99.9% uptime",
		"28d",
		0.999,
		`GraphQL_Errors_total{environment="production",job="unobravo-backend",operationName="getConversations",httpStatusCode=~"5.."}`,
		`GraphQL_Requests_total{environment="production",job="unobravo-backend",operationName="getConversations"}`,
	)

	getConversationsAvailabilityDashboardJSON, err := getConversationsAvailabilitySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating availability dashboard: %v", err)
	}

	outputFileGetConversationsAvailability := "output/get-conversations-availability-slo-dashboard.json"
	if err := os.WriteFile(outputFileGetConversationsAvailability, []byte(getConversationsAvailabilityDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing availability dashboard file: %v", err)
	}

	getConversationsLatencySLO := slo.NewLatencySLO(
		"get-conversations-latency-slo",
		"Get Conversations Latency SLO - 95% requests < 250ms over 28 days",
		"Dashboard to track the monthly latency of the Get Conversations service: 95% of requests should have latency < 250ms",
		"28d",
		0.95,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_bucket{environment="production", operationName="getConversations", job="unobravo-backend", le="250"}`,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_count{environment="production", operationName="getConversations", job="unobravo-backend"}`,
	)

	getConversationsLatencyDashboardJSON, err := getConversationsLatencySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating latency dashboard: %v", err)
	}

	outputFileGetConversationsLatency := "output/get-conversations-latency-slo-dashboard.json"
	if err := os.WriteFile(outputFileGetConversationsLatency, []byte(getConversationsLatencyDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing latency dashboard file: %v", err)
	}

	getMessagesV2AvailabilitySLO := slo.NewAvailabilitySLO(
		"get-messages-v2-availability-slo",
		"Get Messages V2 Availability SLO - 99.9% uptime over 28 days",
		"Dashboard to track the monthly availability of the Get Messages V2 service: 99.9% uptime",
		"28d",
		0.999,
		`GraphQL_Errors_total{environment="production",job="unobravo-backend",operationName="getMessagesV2",httpStatusCode=~"5.."}`,
		`GraphQL_Requests_total{environment="production",job="unobravo-backend",operationName="getMessagesV2"}`,
	)

	getMessagesV2AvailabilityDashboardJSON, err := getMessagesV2AvailabilitySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating availability dashboard: %v", err)
	}

	outputFileGetMessagesV2Availability := "output/get-messages-v2-availability-slo-dashboard.json"
	if err := os.WriteFile(outputFileGetMessagesV2Availability, []byte(getMessagesV2AvailabilityDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing availability dashboard file: %v", err)
	}

	getMessagesV2LatencySLO := slo.NewLatencySLO(
		"get-messages-v2-latency-slo",
		"Get Messages V2 Latency SLO - 95% requests < 250ms over 28 days",
		"Dashboard to track the monthly latency of the Get Messages V2 service: 95% of requests should have latency < 250ms",
		"28d",
		0.95,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_bucket{environment="production", operationName="getMessagesV2", job="unobravo-backend", le="250"}`,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_count{environment="production", operationName="getMessagesV2", job="unobravo-backend"}`,
	)

	getMessagesV2LatencyDashboardJSON, err := getMessagesV2LatencySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating latency dashboard: %v", err)
	}

	outputFileGetMessagesV2Latency := "output/get-messages-v2-latency-slo-dashboard.json"
	if err := os.WriteFile(outputFileGetMessagesV2Latency, []byte(getMessagesV2LatencyDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing latency dashboard file: %v", err)
	}

	sendMessageAvailabilitySLO := slo.NewAvailabilitySLO(
		"send-message-availability-slo",
		"Send Message Availability SLO - 99.9% uptime over 28 days",
		"Dashboard to track the monthly availability of the Send Message service: 99.9% uptime",
		"28d",
		0.999,
		`GraphQL_Errors_total{environment="production",job="unobravo-backend",operationName="sendMessage",httpStatusCode=~"5.."}`,
		`GraphQL_Requests_total{environment="production",job="unobravo-backend",operationName="sendMessage"}`,
	)

	sendMessageAvailabilityDashboardJSON, err := sendMessageAvailabilitySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating availability dashboard: %v", err)
	}

	outputFileSendMessageAvailability := "output/send-message-availability-slo-dashboard.json"
	if err := os.WriteFile(outputFileSendMessageAvailability, []byte(sendMessageAvailabilityDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing availability dashboard file: %v", err)
	}

	sendMessageLatencySLO := slo.NewLatencySLO(
		"send-message-latency-slo",
		"Send Message Latency SLO - 95% requests < 400ms over 28 days",
		"Dashboard to track the monthly latency of the Send Message service: 95% of requests should have latency < 400ms",
		"28d",
		0.95,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_bucket{environment="production", operationName="sendMessage", job="unobravo-backend", le="400"}`,
		`GraphQL_WebTransactionTimeHistogram_milliseconds_count{environment="production", operationName="sendMessage", job="unobravo-backend"}`,
	)

	sendMessageLatencyDashboardJSON, err := sendMessageLatencySLO.BuildJSON()
	if err != nil {
		log.Fatalf("Error generating latency dashboard: %v", err)
	}

	outputFileSendMessageLatency := "output/send-message-latency-slo-dashboard.json"
	if err := os.WriteFile(outputFileSendMessageLatency, []byte(sendMessageLatencyDashboardJSON), 0644); err != nil {
		log.Fatalf("Error writing latency dashboard file: %v", err)
	}

}
