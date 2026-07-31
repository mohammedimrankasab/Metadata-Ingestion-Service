// Package processor handles metadata processing,
// retry logic and persistence to the configured sink.
package processor

import (
	"context"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/metrics"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/retry"
	inSink "github.com/mohammedimrankasab/metadata-ingestion-service/internal/sink"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Processor struct {
	logger *zap.Logger
	sink   inSink.Sink
}

func NewProcessor(
	logger *zap.Logger,
	sink inSink.Sink,
) *Processor {

	return &Processor{
		logger: logger,
		sink:   sink,
	}
}

var retryConfig = retry.Config{
	MaxRetries: 3,
	BaseDelay:  500 * time.Millisecond,
}

func (p *Processor) Process(
	ctx context.Context,
	job models.MetadataJob,
) error {

	tracer := otel.Tracer(
		"processor",
	)

	ctx, span := tracer.Start(
		ctx,
		"ProcessMetadata",
	)

	defer span.End()

	span.SetAttributes(
		attribute.String(
			"job.id",
			job.ID,
		),
		attribute.String(
			"connector",
			job.Connector,
		),
	)

	start := time.Now()

	defer func() {

		metrics.ProcessingDuration.Observe(
			time.Since(start).Seconds(),
		)

	}()

	p.logger.Debug(
		"processing metadata",
		zap.String(
			"job_id",
			job.ID,
		),
		zap.String(
			"connector",
			job.Connector,
		),
	)

	err := retry.Do(
		ctx,
		retryConfig,
		func() error {

			sinkStart := time.Now()

			err := p.sink.Write(
				ctx,
				job.Metadata,
			)

			metrics.SinkProcessingDuration.Observe(
				time.Since(sinkStart).Seconds(),
			)

			return err
		},
	)

	if err != nil {

		span.RecordError(err)

		span.SetStatus(
			codes.Error,
			err.Error(),
		)

		p.logger.Error(
			"metadata processing failed",
			zap.String(
				"job_id",
				job.ID,
			),
			zap.String(
				"connector",
				job.Connector,
			),
			zap.Error(err),
		)

		metrics.JobsFailed.Inc()

		return err
	}

	metrics.JobsProcessed.
		WithLabelValues(
			job.Connector,
		).
		Inc()

	p.logger.Debug(
		"metadata processed",
		zap.String(
			"job_id",
			job.ID,
		),
	)

	return nil
}
