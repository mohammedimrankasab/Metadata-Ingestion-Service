package ingestion

import (
	"context"
	"sync"

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

	const workerCount = 4

	jobsCh := make(chan models.MetadataJob, 100)

	var wg sync.WaitGroup

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)

		go StartWorker(
			ctx,
			i,
			s.logger,
			&wg,
			jobsCh,
			s.processor,
		)
	}

	for _, connector := range s.connectors {

		metadataList, err := connector.FetchMetadata(ctx, nil)
		if err != nil {
			close(jobsCh)
			wg.Wait()
			return err
		}

		for _, metadata := range metadataList {
			job := models.NewJob(
				connector.Name(),
				metadata,
			)

			jobsCh <- job
		}
	}

	close(jobsCh)

	wg.Wait()

	s.logger.Info("Metadata ingestion completed")

	return nil

}
