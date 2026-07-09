package main

import (
	"context"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/ingestion"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/logger"
)

func main() {

	ctx := context.Background()

	service := ingestion.New(
		connectors.NewPowerBIConnector(),
	)

	if err := service.Run(ctx); err != nil {
		logger.Error("Error running ingestion service: " + err.Error())
		panic(err)
	}
}
