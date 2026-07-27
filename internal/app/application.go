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
	"go.uber.org/zap"
)

type Application struct {
	Logger           *zap.Logger
	Server           *server.Server
	IngestionService *ingestion.Service
	Config           *config.Config
}

func NewApplication(
	cfg *config.Config,
	logger *zap.Logger,
) (*Application, error) {

	powerBI := connectors.NewPowerBIConnector(logger)
	consoleSink := inSink.NewConsoleSink(logger)

	metadataProcessor := processor.NewProcessor(
		logger,
		consoleSink,
	)
	service := ingestion.New(
		logger,
		cfg,
		metadataProcessor,
		powerBI,
	)
	httpServer := server.New(logger, cfg, service)

	return &Application{
		Logger:           logger,
		Server:           httpServer,
		IngestionService: service,
		Config:           cfg,
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

	return app.Server.Shutdown(shutdownCtx)
}
