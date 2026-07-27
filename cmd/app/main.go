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

	logger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	shutdownTracer, err := telemetry.InitTracer()
	if err != nil {
		logger.Fatal(
			"failed to initialize telemetry",
			zap.Error(err),
		)
	}

	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			logger.Error(
				"failed to shutdown tracer",
				zap.Error(err),
			)
		}
	}()

	metrics.Register()

	application, err := app.NewApplication(cfg, logger)
	if err != nil {
		logger.Fatal(
			"unable to create application",
			zap.Error(err),
		)
	}

	if err := application.Run(ctx); err != nil {
		logger.Fatal(
			"application failed",
			zap.Error(err),
		)
	}
}
