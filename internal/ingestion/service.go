package ingestion

import (
	"sync"
	"time"

	"github.com/mohammedimrankasab/go-concurrency-examples/internal/logger"
	"github.com/mohammedimrankasab/go-concurrency-examples/internal/models"
)

type IngestionService struct{}

func NewIngestionService() *IngestionService {
	return &IngestionService{}
}

func (s *IngestionService) IngestData(workspaces []string) error {

	var wg sync.WaitGroup

	start := time.Now()

	for index, workspace := range workspaces {
		wg.Add(1)
		go func(id int, w string) {
			defer wg.Done()
			// Implementation for ingesting workspace
			// Simulating ingestion with sleep

			metadata := models.Metadata{
				ID:         id,
				Source:     "PowerBI",
				Name:       w + "-Dashboard",
				LastUpdate: time.Now().Format(time.RFC3339),
			}
			time.Sleep(2 * time.Second)
			logger.Info("completed ingestion for " + metadata.Name)
		}(index+1, workspace)
	}
	wg.Wait()
	logger.Info("Total time taken for ingestion: " + time.Since(start).String())
	return nil
}
