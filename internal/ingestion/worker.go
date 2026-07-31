// Package ingestion orchestrates metadata collection and
// concurrent processing using a configurable worker pool.
package ingestion

import (
	"context"
	"sync"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

type JobProcessor interface {
	Process(context.Context, models.MetadataJob) error
}

func worker(
	ctx context.Context,
	id int,
	logger *zap.Logger,
	wg *sync.WaitGroup,
	jobs <-chan models.MetadataJob,
	processor JobProcessor,
) {

	defer func() {
		logger.Debug(
			"worker exited",
			zap.Int("worker", id),
		)

		wg.Done()
	}()

	logger.Debug(
		"worker started",
		zap.Int("worker", id),
	)

	for {

		select {

		case <-ctx.Done():

			logger.Debug(
				"worker stopped",
				zap.Int("worker", id),
			)

			return

		case job, ok := <-jobs:

			if !ok {

				logger.Debug(
					"job channel closed",
					zap.Int("worker", id),
				)

				return
			}

			err := processor.Process(
				ctx,
				job,
			)

			if err != nil {

				logger.Error(
					"failed processing metadata",
					zap.Int("worker", id),
					zap.String("job_id", job.ID),
					zap.String("connector", job.Connector),
					zap.Error(err),
				)
			}

			close(job.Done)
		}
	}
}
