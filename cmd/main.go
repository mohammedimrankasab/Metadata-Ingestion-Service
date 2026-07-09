package main

import "github.com/mohammedimrankasab/go-concurrency-examples/internal/ingestion"

func main() {

	workspaces := []string{
		"Marketing",
		"Sales",
		"Finance",
		"HR",
		"Engineering",
	}

	service := ingestion.NewIngestionService()
	_ = service.IngestData(workspaces)
}
