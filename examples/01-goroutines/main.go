package main

import (
	"sync"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/logger"
)

func ingestWorkspace(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Implementation for ingesting workspace

	logger.Log.Info("Ingesting workspace: " + name)

	time.Sleep(2 * time.Second)

	logger.Log.Info("Finished ingesting workspace: " + name)
}

func main() {

	workspaces := []string{
		"Finance",
		"Sales",
		"Marketing",
		"HR",
		"Engineering",
	}

	var wg sync.WaitGroup
	start := time.Now()
	logger.Log.Info("Starting the ingestion at:" + start.String())
	for _, workspace := range workspaces {
		wg.Add(1)
		go ingestWorkspace(workspace, &wg)
	}

	wg.Wait()
	logger.Log.Info("All workspaces ingested successfully.")
	logger.Log.Info("Total time taken:" + time.Since(start).String())
}
