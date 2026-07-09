package ingestion

import (
	"context"
	"sync"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/processor"
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

	const workerCount = 4

	jobCh := make(chan models.MetadataJob, 100)

	var wg sync.WaitGroup
	p := processor.NewProcessor()

	for i := 0; i <= workerCount; i++ {
		wg.Add(1)
		go StartWorker(ctx, i, &wg, jobCh, p)
	}
	for _, connector := range s.Connector {
		metadata, err := connector.FetchMetadata(ctx, nil)
		if err != nil {
			return err
		}

		for _, item := range metadata {
			jobCh <- models.NewJob(item)
		}
	}
	close(jobCh)
	wg.Wait()

	return nil

}
