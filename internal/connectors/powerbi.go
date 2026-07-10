package connectors

import (
	"context"
	"fmt"
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
	select {

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-time.After(10 * time.Second):

	}
	metadataList := make([]models.Metadata, 0, 200)

	for i := 1; i <= 200; i++ {

		metadataList = append(metadataList,
			models.NewMetadata(
				uuid.NewString(),
				fmt.Sprintf("Dashboard-%d", i),
				models.DashboardType,
				"Finance",
				p.Name(),
				time.Now(),
			),
		)

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
