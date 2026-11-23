package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EventsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pulsestream_events_received_total",
		Help: "Amount of pulse stream events received by type",
	}, []string{"type"})

	PublishErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pulsestream_publish_errors_total",
		Help: "Errors encountered while publishing events",
	})

	PublishDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "pulsestream_publish_duration_seconds",
		Help:    "Time to process events",
		Buckets: prometheus.DefBuckets,
	})
)
