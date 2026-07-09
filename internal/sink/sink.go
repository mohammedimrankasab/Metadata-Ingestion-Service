package sink

import (
	"context"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
)

type Sink interface {
	Write(
		ctx context.Context,
		metadata models.Metadata,
	) error
}
