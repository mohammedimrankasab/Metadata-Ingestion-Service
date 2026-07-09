package main

import (
	"context"
	"log"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/app"
	"go.uber.org/zap"
)

func main() {

	application, err := app.NewApplication()

	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if err := application.Run(ctx); err != nil {

		application.Components.Logger.Fatal(
			"application failed",
			zap.Error(err),
		)

	}

}
