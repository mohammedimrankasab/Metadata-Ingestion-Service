package sink

import (
	"context"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

type ConsoleSink struct {
	logger *zap.Logger
}

func NewConsoleSink(
	logger *zap.Logger,
) *ConsoleSink {
	return &ConsoleSink{
		logger: logger,
	}
}

func (c *ConsoleSink) Write(
	ctx context.Context,
	metadata models.Metadata,
) error {

	c.logger.Info(
		"Metadata Written",
		zap.String("connector", metadata.Source),
		zap.String("workspace", metadata.Workspace),
		zap.String("name", metadata.Name),
		zap.String("type", string(metadata.Type)),
	)

	return nil
}
