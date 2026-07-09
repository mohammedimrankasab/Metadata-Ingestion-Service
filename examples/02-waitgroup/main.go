package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/mohammedimrankasab/go-concurrency-examples/internal/logger"
)

func ingest(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Implementation for ingesting data
	strID := strconv.Itoa(id)
	logger.Info("Processing metadata for id: " + strID)
	time.Sleep(time.Second)
	logger.Info("Completed metadata ingestion for id: " + strID)
}

func main() {
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go ingest(i, &wg)

	}
	wg.Wait()
	logger.Info("All metadata ingested successfully.")
}
