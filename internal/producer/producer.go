package producer

import (
	"context"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

type Producer struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Producer {
	return &Producer{
		logger: logger,
	}
}

func (p *Producer) Produce(
	ctx context.Context,
	metadata []models.Metadata,
	jobs chan<- models.MetadataJob,
) error {

	for _, item := range metadata {
		job := models.NewJob(item)
		select {
		case <-ctx.Done():
			return ctx.Err()

		case jobs <- job:
			p.logger.Debug(
				"Job submitted",
				zap.String("jobId", job.ID),
			)
		}
	}
	return nil
}
