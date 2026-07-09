package main

import (
	"sync"
	"time"

	"github.com/mohammedimrankasab/go-concurrency-examples/internal/logger"
)

func ingestWorkspace(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Implementation for ingesting workspace

	logger.Info("Ingesting workspace: " + name)

	time.Sleep(2 * time.Second)

	logger.Info("Finished ingesting workspace: " + name)
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
	logger.Info("Starting the ingestion at:" + start.String())
	for _, workspace := range workspaces {
		wg.Add(1)
		go ingestWorkspace(workspace, &wg)
	}

	wg.Wait()
	logger.Info("All workspaces ingested successfully.")
	logger.Info("Total time taken:" + time.Since(start).String())
}
