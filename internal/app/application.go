package app

import (
	"context"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/ingestion"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/processor"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/server"
	inSink "github.com/mohammedimrankasab/metadata-ingestion-service/internal/sink"
	"go.uber.org/zap"
)

type Application struct {
	Components       *Components
	IngestionService *ingestion.Service
	Config           config.Config
}

func NewApplication(
	cfg config.Config,
	logger *zap.Logger,
) (*Application, error) {

	powerBI := connectors.NewPowerBIConnector(logger)
	consoleSink := inSink.NewConsoleSink(logger)

	processor := processor.NewProcessor(
		logger,
		consoleSink,
	)
	service := ingestion.New(
		logger,
		cfg,
		processor,
		powerBI,
	)
	metricsServer := server.New(logger, cfg, service)
	components := &Components{
		Logger: logger,
		Server: metricsServer,
	}
	return &Application{
		Components:       components,
		IngestionService: service,
		Config:           cfg,
	}, nil
}

type Components struct {
	Logger *zap.Logger
	Server *server.Server
}

func (app *Application) Run(ctx context.Context) error {
	app.Components.Logger.Info("Application starting")
	app.Components.Server.Start()
	defer func() {
		_ = app.Components.Logger.Sync()
	}()
	if err := app.IngestionService.Run(ctx); err != nil {
		return err
	}

	app.Components.Logger.Info("Application stopped")

	return nil
}
