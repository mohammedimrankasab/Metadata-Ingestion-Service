// Package ingestion orchestrates metadata collection and
// concurrent processing using a configurable worker pool.
package ingestion

import (
	"context"
	"errors"
	"sync"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/processor"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Service struct {
	logger     *zap.Logger
	config     *config.Config
	processor  *processor.Processor
	connectors []connectors.Connector

	workerCtx    context.Context
	workerCancel context.CancelFunc

	jobs chan models.MetadataJob
	wg   sync.WaitGroup

	startOnce sync.Once
}

func New(
	logger *zap.Logger,
	config *config.Config,
	processor *processor.Processor,
	connectors ...connectors.Connector,
) *Service {

	workerCtx, cancel := context.WithCancel(context.Background())

	return &Service{
		logger:       logger,
		config:       config,
		processor:    processor,
		connectors:   connectors,
		workerCtx:    workerCtx,
		workerCancel: cancel,
		jobs: make(
			chan models.MetadataJob,
			config.JobQueueSize,
		),
	}
}

func (s *Service) StartWorkers() {

	s.startOnce.Do(func() {

		s.logger.Info(
			"starting worker pool",
			zap.Int(
				"workers",
				s.config.WorkerCount,
			),
			zap.Int(
				"queue_size",
				s.config.JobQueueSize,
			),
		)

		for i := 1; i <= s.config.WorkerCount; i++ {

			s.wg.Add(1)

			go worker(
				s.workerCtx,
				i,
				s.logger,
				&s.wg,
				s.jobs,
				s.processor,
			)
		}
	})
}

func (s *Service) Run(ctx context.Context) error {

	tracer := otel.Tracer(
		"ingestion",
	)

	ctx, span := tracer.Start(
		ctx,
		"metadata ingestion",
	)

	defer span.End()

	if len(s.connectors) == 0 {

		err := errors.New(
			"no connectors configured",
		)

		span.RecordError(err)

		span.SetStatus(
			codes.Error,
			err.Error(),
		)

		return err
	}

	var pending []<-chan struct{}

	for _, connector := range s.connectors {

		connectorCtx, connectorSpan := tracer.Start(
			ctx,
			"fetch metadata",
		)

		connectorSpan.SetAttributes(
			attribute.String(
				"connector",
				connector.Name(),
			),
		)

		s.logger.Info(
			"processing connector",
			zap.String(
				"connector",
				connector.Name(),
			),
		)

		metadataList, err := connector.FetchMetadata(
			connectorCtx,
			nil,
		)

		if err != nil {

			connectorSpan.RecordError(err)

			connectorSpan.SetStatus(
				codes.Error,
				err.Error(),
			)

			connectorSpan.End()

			return err
		}

		connectorSpan.SetAttributes(
			attribute.Int(
				"metadata.count",
				len(metadataList),
			),
		)

		connectorSpan.End()

		s.logger.Info(
			"metadata fetched",
			zap.String(
				"connector",
				connector.Name(),
			),
			zap.Int(
				"count",
				len(metadataList),
			),
		)

		for _, metadata := range metadataList {

			job := models.NewJob(
				connector.Name(),
				metadata,
			)

			job.Done = make(chan struct{})

			select {

			case <-ctx.Done():

				return ctx.Err()

			case s.jobs <- job:
				pending = append(
					pending,
					job.Done,
				)
			}
		}
	}
	for _, done := range pending {
		<-done
	}
	s.logger.Info(
		"metadata ingestion completed",
	)

	return ctx.Err()
}

func (s *Service) Shutdown() {

	s.logger.Info(
		"stopping ingestion workers",
	)

	s.workerCancel()

	close(s.jobs)

	s.wg.Wait()

	s.logger.Info(
		"all workers stopped",
	)
}
