package connectors

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

type PowerBIConnector struct {
	logger *zap.Logger
}

func NewPowerBIConnector(logger *zap.Logger) *PowerBIConnector {
	return &PowerBIConnector{
		logger: logger,
	}
}

func (p *PowerBIConnector) Name() string {
	return "PowerBI"
}

func (p *PowerBIConnector) FetchMetadata(ctx context.Context, lastSyncTime *time.Time) ([]models.Metadata, error) {

	p.logger.Info("Fetching metadata from PowerBI connector...")

	metadata := []models.Metadata{
		{
			ID:           uuid.NewString(),
			Name:         "Finance Dashboard",
			Source:       p.Name(),
			Type:         models.DashboardType,
			Workspace:    "Finance",
			LastModified: time.Now(),
		},
		{
			ID:           uuid.NewString(),
			Name:         "Sales Report",
			Source:       p.Name(),
			Type:         models.ReportType,
			Workspace:    "Sales",
			LastModified: time.Now(),
		},
	}

	if lastSyncTime != nil {
		filteredMetadata := make([]models.Metadata, 0)

		for _, m := range metadata {
			if m.LastModified.After(*lastSyncTime) {
				filteredMetadata = append(filteredMetadata, m)
			}
		}

		return filteredMetadata, nil
	}
	return metadata, nil
}
