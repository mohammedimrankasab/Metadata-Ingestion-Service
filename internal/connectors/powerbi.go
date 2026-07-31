package connectors

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

const sampleMetadataCount = 200

type PowerBIConnector struct {
	BaseConnector
}

func NewPowerBIConnector(
	logger *zap.Logger,
) *PowerBIConnector {

	return &PowerBIConnector{
		BaseConnector: NewBaseConnector(logger),
	}
}

func (p *PowerBIConnector) Name() string {
	return "PowerBI"
}

func (p *PowerBIConnector) FetchMetadata(
	ctx context.Context,
	lastSyncTime *time.Time,
) ([]models.Metadata, error) {

	ctx, endTrace := p.StartTrace(
		ctx,
		"PowerBI.FetchMetadata",
	)

	defer endTrace()

	p.logger.Info(
		"fetching metadata from PowerBI",
	)

	select {

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-time.After(
		500 * time.Millisecond,
	):

	}

	now := time.Now()

	metadata := make(
		[]models.Metadata,
		0,
		sampleMetadataCount,
	)

	for i := 1; i <= sampleMetadataCount; i++ {

		metadata = append(
			metadata,
			models.NewMetadata(
				uuid.NewString(),
				fmt.Sprintf(
					"Dashboard-%d",
					i,
				),
				models.DashboardType,
				"Finance",
				p.Name(),
				now,
			),
		)
	}

	metadata = p.FilterByLastSync(
		metadata,
		lastSyncTime,
	)

	p.logger.Info(
		"metadata fetched",
		zap.String(
			"connector",
			p.Name(),
		),
		zap.Int(
			"count",
			len(metadata),
		),
	)

	return metadata, nil
}

func (p *PowerBIConnector) Health(
	ctx context.Context,
) error {

	p.logger.Debug(
		"checking connector health",
	)

	return nil
}
