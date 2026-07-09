package ingestion

import (
	"context"
	"sync"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/processor"
	"go.uber.org/zap"
)

func StartWorker(
	ctx context.Context,
	id int,
	logger *zap.Logger,
	wg *sync.WaitGroup,
	jobs <-chan models.MetadataJob,
	processor *processor.Processor,
) {

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info(
				"Worker stopped",
				zap.Int("worker", id),
			)
			return
		case job, ok := <-jobs:
			if !ok {

				logger.Info(
					"Job channel closed",
					zap.Int("worker", id),
				)

				return
			}
			if err := processor.Process(ctx, job); err != nil {
				logger.Error(
					"Failed processing",
					zap.Error(err),
				)
			}
		}
	}
}
