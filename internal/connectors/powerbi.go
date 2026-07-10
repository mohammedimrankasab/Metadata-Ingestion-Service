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

	m1 := models.NewMetadata(
		uuid.NewString(),
		"Finance Dashboard",
		models.DashboardType,
		"Finance",
		p.Name(),
		time.Now(),
	)
	m2 := models.NewMetadata(
		uuid.NewString(),
		"Sales Report",
		models.ReportType,
		"Sales",
		p.Name(),
		time.Now(),
	)
	metadataList := []models.Metadata{
		m1,
		m2,
	}

	if lastSyncTime == nil {
		return metadataList, nil
	}

	filtered := make([]models.Metadata, 0, len(metadataList))

	for _, m := range metadataList {
		if m.LastModified.After(*lastSyncTime) {
			filtered = append(filtered, m)
		}
	}

	return filtered, nil
}
func (p *PowerBIConnector) Health(ctx context.Context) error {
	return nil
}
