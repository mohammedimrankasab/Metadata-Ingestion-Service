package main

import (
	"context"
	"log"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/ingestion"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/logger"
)

func main() {

	if err := logger.Init(); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	ctx := context.Background()

	service := ingestion.New(
		connectors.NewPowerBIConnector(),
	)

	if err := service.Run(ctx); err != nil {
		logger.Log.Fatal("ingestion failed")
	}
}
