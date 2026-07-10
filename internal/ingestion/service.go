package ingestion

import (
	"context"
	"runtime"
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
}

func New(logger *zap.Logger, processor *processor.Processor, connectors ...connectors.Connector) *Service {
	return &Service{
		logger:     logger,
		processor:  processor,
		connectors: connectors,
	}
}

func (s *Service) Run(ctx context.Context) error {

	cfg := inConfig.Config{
		WorkerCount:  runtime.NumCPU(),
		JobQueueSize: 100,
	}

	jobs := make(chan models.MetadataJob, cfg.JobQueueSize)

	var wg sync.WaitGroup

	for i := 1; i <= cfg.WorkerCount; i++ {
		wg.Add(1)

		go StartWorker(
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
		wg.Wait()
	}()
	for _, connector := range s.connectors {

		metadataList, err := connector.FetchMetadata(ctx, nil)
		if err != nil {
			return err
		}

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
