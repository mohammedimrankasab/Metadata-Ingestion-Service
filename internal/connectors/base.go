// Package connectors provides metadata source connectors.
package connectors

import (
	"context"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type BaseConnector struct {
	logger *zap.Logger
}

func NewBaseConnector(
	logger *zap.Logger,
) BaseConnector {
	return BaseConnector{
		logger: logger,
	}
}

func (b BaseConnector) FilterByLastSync(
	metadata []models.Metadata,
	lastSyncTime *time.Time,
) []models.Metadata {

	if lastSyncTime == nil {
		return metadata
	}

	filtered := make(
		[]models.Metadata,
		0,
		len(metadata),
	)

	for _, item := range metadata {

		if item.LastModified.After(*lastSyncTime) {
			filtered = append(
				filtered,
				item,
			)
		}
	}

	return filtered
}

func (b BaseConnector) StartTrace(
	ctx context.Context,
	name string,
) (
	context.Context,
	func(),
) {

	tracer := otel.Tracer(
		"metadata-ingestion-service/connectors",
	)

	ctx, span := tracer.Start(
		ctx,
		name,
	)

	return ctx, func() {
		span.End()
	}
}
