package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var MetricRequestsTotal = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "shorter",
		Subsystem:  "http",
		Name:       "requests",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.01},
	},
	[]string{"status"},
)
