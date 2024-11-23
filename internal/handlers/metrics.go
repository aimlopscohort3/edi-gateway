package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Define Prometheus metrics
var (
	InboundCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "inbound_requests_total",
		Help: "Total number of inbound EDI transactions.",
	})
	OutboundCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "outbound_requests_total",
		Help: "Total number of outbound EDI transactions.",
	})
)

// RegisterMetrics registers the Prometheus metrics
func RegisterMetrics() {
	prometheus.MustRegister(InboundCounter, OutboundCounter)
}