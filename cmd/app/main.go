package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/app"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/metrics"
	"go.uber.org/zap"
)

func main() {
	metrics.Register()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	application, err := app.NewApplication()

	if err != nil {
		log.Fatal(err)
	}
	if err := application.Run(ctx); err != nil {

		application.Components.Logger.Fatal(
			"application failed",
			zap.Error(err),
		)

	}

}
