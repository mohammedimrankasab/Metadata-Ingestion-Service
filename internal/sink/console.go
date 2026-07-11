package sink

import (
	"context"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.opentelemetry.io/otel"
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
	tracer := otel.Tracer("sink")

	_, span := tracer.Start(
		ctx,
		"WriteMetadata",
	)

	defer span.End()
	return nil
}
