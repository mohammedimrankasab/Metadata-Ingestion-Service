package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/app"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/logger"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/metrics"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/telemetry"
	"go.uber.org/zap"
)

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := config.Load()

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	metrics.Register()

	application, err := app.NewApplication(cfg, logger)

	if err != nil {
		logger.Fatal(
			"unable to create application",
			zap.Error(err),
		)
	}
	if err := application.Run(ctx); err != nil {

		application.Components.Logger.Fatal(
			"application failed",
			zap.Error(err),
		)

	}
	shutdownTracer, err := telemetry.Init()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			log.Printf("error shutting down tracer: %v", err)
		}
	}()

}
