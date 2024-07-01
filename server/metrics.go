package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func record_metrics(method string, process func()) {
	start := time.Now()
	process()
	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues(method).Observe(duration)
	requestCounter.WithLabelValues(method).Inc()
}

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "method_cost",
			Help:    "method cost",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_counter",
			Help: "request counter",
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestCounter)
}
