package processor

import (
	"context"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/logger"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Process(ctx context.Context, job models.MetadataJob) error {
	logger.Log.Info("Processing metadata",
		zap.String("name", job.Metadata.Name),
		zap.String("workspace", job.Metadata.Workspace),
	)
	time.Sleep(500 * time.Millisecond)
	return nil
}
