package ingestion

import (
	"context"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
)

type Service struct {
	Connector []connectors.Connector
}

func New(connectors ...connectors.Connector) *Service {
	return &Service{
		Connector: connectors,
	}
}

func (s *Service) Run(ctx context.Context) error {

	for _, connector := range s.Connector {
		_, err := connector.FetchMetadata(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil

}
