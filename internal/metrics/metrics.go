package metrics

import "github.com/prometheus/client_golang/prometheus"

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
			Help: "Total failed jobs",
		},
	)

	RetryCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "metadata_job_retries_total",
			Help: "Retry count",
		},
	)

	ProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "metadata_processing_duration_seconds",
			Help: "Processing duration",
		},
	)
)

func Register() {

	prometheus.MustRegister(
		JobsProcessed,
		JobsFailed,
		RetryCount,
		ProcessingDuration,
	)

}
