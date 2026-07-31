package app

import (
	"context"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/ingestion"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/processor"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/server"
	inSink "github.com/mohammedimrankasab/metadata-ingestion-service/internal/sink"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/telemetry"
	"go.uber.org/zap"
)

type Application struct {
	ctx              context.Context
	Logger           *zap.Logger
	Server           *server.Server
	IngestionService *ingestion.Service
	Config           *config.Config
	tracerShutdown   func(context.Context) error
}

func NewApplication(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.Logger,
) (*Application, error) {

	tracerShutdown, err := telemetry.InitTracer(cfg.EnableTracing)
	if err != nil {
		return nil, err
	}

	powerBI := connectors.NewPowerBIConnector(logger)
	csv := connectors.NewCSVConnector(logger)

	selectedSink, err := inSink.New(cfg, logger)
	if err != nil {
		return nil, err
	}

	metadataProcessor := processor.NewProcessor(
		logger,
		selectedSink,
	)
	ingestionService := ingestion.New(
		logger,
		cfg,
		metadataProcessor,
		powerBI,
		csv,
	)
	ingestionService.StartWorkers()
	httpServer := server.New(ctx, logger, cfg, ingestionService)

	return &Application{
		Logger:           logger,
		Server:           httpServer,
		IngestionService: ingestionService,
		Config:           cfg,
		tracerShutdown:   tracerShutdown,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	app.Logger.Info("Application starting")

	defer func() {
		_ = app.Logger.Sync()
	}()

	app.Server.Start()
	app.Logger.Info("application started")

	<-ctx.Done()

	app.Logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err := app.Server.Shutdown(shutdownCtx)

	app.IngestionService.Shutdown()

	if app.tracerShutdown != nil {
		if tracerErr := app.tracerShutdown(
			shutdownCtx,
		); tracerErr != nil {
			app.Logger.Error(
				"failed to shutdown tracer",
				zap.Error(tracerErr),
			)
		}
	}

	return err
}
