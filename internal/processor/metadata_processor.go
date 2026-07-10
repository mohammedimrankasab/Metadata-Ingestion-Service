package processor

import (
	"context"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	inSink "github.com/mohammedimrankasab/metadata-ingestion-service/internal/sink"
	"go.uber.org/zap"
)

type Processor struct {
	logger *zap.Logger
	sink   inSink.Sink
}

func NewProcessor(
	logger *zap.Logger,
	sink inSink.Sink,
) *Processor {
	return &Processor{
		logger: logger,
		sink:   sink,
	}
}

func (p *Processor) Process(ctx context.Context, job models.MetadataJob) error {

	p.logger.Debug(
		"Processing metadata",
		zap.String("jobId", job.ID),
		zap.String("connector", job.Connector),
	)

	return p.sink.Write(
		ctx,
		job.Metadata,
	)
}
