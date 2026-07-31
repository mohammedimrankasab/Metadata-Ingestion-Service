package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	JobsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "metadata_jobs_processed_total",
			Help: "Total processed metadata jobs",
		},
		[]string{"connector"},
	)

	JobsFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "metadata_jobs_failed_total",
			Help: "Total number of failed metadata jobs.",
		},
	)

	RetryCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "metadata_job_retries_total",
			Help: "Total number of retry attempts.",
		},
	)

	ProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "metadata_processing_duration_seconds",
			Help:    "Time taken to process a metadata job.",
			Buckets: prometheus.DefBuckets,
		},
	)

	SinkProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "metadata_sink_processing_duration_seconds",
			Help: "Time taken by sink write operation",
		},
	)
)
var registerOnce sync.Once

func Register() {
	registerOnce.Do(func() {
		prometheus.MustRegister(
			JobsProcessed,
			JobsFailed,
			RetryCount,
			ProcessingDuration,
		)
	})
}
