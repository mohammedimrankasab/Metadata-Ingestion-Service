package connectors

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

type CSVConnector struct {
	BaseConnector
}

func NewCSVConnector(
	logger *zap.Logger,
) *CSVConnector {

	return &CSVConnector{
		BaseConnector: NewBaseConnector(logger),
	}
}

func (c *CSVConnector) Name() string {
	return "CSV"
}

func (c *CSVConnector) FetchMetadata(
	ctx context.Context,
	lastSyncTime *time.Time,
) ([]models.Metadata, error) {

	ctx, endTrace := c.StartTrace(
		ctx,
		"CSV.FetchMetadata",
	)

	defer endTrace()

	now := time.Now()

	metadata := []models.Metadata{

		models.NewMetadata(
			uuid.NewString(),
			"employees.csv",
			models.TableType,
			"Local",
			c.Name(),
			now,
		),
	}

	metadata = c.FilterByLastSync(
		metadata,
		lastSyncTime,
	)

	c.logger.Info(
		"metadata fetched",
		zap.String(
			"connector",
			c.Name(),
		),
		zap.Int(
			"count",
			len(metadata),
		),
	)

	return metadata, nil
}

func (c *CSVConnector) Health(
	ctx context.Context,
) error {

	c.logger.Debug(
		"checking connector health",
	)

	return nil
}
