package connectors

import (
	"context"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
)

type Connector interface {
	Name() string
	FetchMetadata(ctx context.Context, lastSyncTime *time.Time) ([]models.Metadata, error)
}
