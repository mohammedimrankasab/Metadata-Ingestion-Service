package processor

import (
	"context"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

type Processor struct {
	logger *zap.Logger
}

func NewProcessor(logger *zap.Logger) *Processor {
	return &Processor{
		logger: logger,
	}
}

func (p *Processor) Process(ctx context.Context, job models.MetadataJob) error {
	p.logger.Info("Processing metadata",
		zap.String("name", job.Metadata.Name),
		zap.String("workspace", job.Metadata.Workspace),
	)
	time.Sleep(500 * time.Millisecond)
	return nil
}
