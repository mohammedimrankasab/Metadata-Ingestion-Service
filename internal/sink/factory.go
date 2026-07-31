package sink

import (
	"fmt"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"go.uber.org/zap"
)

func New(
	cfg *config.Config,
	logger *zap.Logger,
) (Sink, error) {

	switch cfg.SinkType {

	case "console":
		return NewConsoleSink(logger), nil

	case "json":
		return nil, fmt.Errorf("json sink not implemented")

	case "opensearch":
		return NewOpenSearchSink(
			logger,
			cfg.OpenSearchURL,
			cfg.OpenSearchIndex,
		), nil

	default:
		return nil, fmt.Errorf(
			"unsupported sink type: %s",
			cfg.SinkType,
		)
	}
}
