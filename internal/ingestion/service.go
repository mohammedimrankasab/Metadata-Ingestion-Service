// Package ingestion orchestrates metadata collection and
// concurrent processing using a configurable worker pool.
package ingestion

import (
	"context"
	"errors"
	"sync"

	inConfig "github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/processor"
	"go.uber.org/zap"
)

type Service struct {
	logger     *zap.Logger
	processor  *processor.Processor
	connectors []connectors.Connector
	config     *inConfig.Config
}

func New(logger *zap.Logger, config *inConfig.Config, processor *processor.Processor, connectors ...connectors.Connector) *Service {
	return &Service{
		logger:     logger,
		config:     config,
		processor:  processor,
		connectors: connectors,
	}
}

func (s *Service) Run(ctx context.Context) error {
	if len(s.connectors) == 0 {
		return errors.New("no connectors configured")
	}

	jobs := make(chan models.MetadataJob, s.config.JobQueueSize)

	var wg sync.WaitGroup

	for i := 1; i <= s.config.WorkerCount; i++ {
		wg.Add(1)

		go worker(
			ctx,
			i,
			s.logger,
			&wg,
			jobs,
			s.processor,
		)
	}
	defer func() {
		close(jobs)

		s.logger.Info("waiting for workers")

		wg.Wait()

		s.logger.Info("all workers completed")
	}()
	for _, connector := range s.connectors {
		s.logger.Info(
			"processing connector",
			zap.String("connector", connector.Name()),
		)
		metadataList, err := connector.FetchMetadata(ctx, nil)
		if err != nil {
			return err
		}
		s.logger.Info(
			"metadata fetched",
			zap.String("connector", connector.Name()),
			zap.Int("count", len(metadataList)),
		)
		for _, metadata := range metadataList {
			job := models.NewJob(
				connector.Name(),
				metadata,
			)

			select {
			case <-ctx.Done():
				s.logger.Info("Stopping job submission")
				return ctx.Err()

			case jobs <- job:
			}
		}
	}
	s.logger.Info("Metadata ingestion completed")

	return nil

}
