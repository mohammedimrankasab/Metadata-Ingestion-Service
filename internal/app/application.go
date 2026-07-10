package app

import (
	"context"

	inConfig "github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/ingestion"
	inLog "github.com/mohammedimrankasab/metadata-ingestion-service/internal/logger"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/processor"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/server"
	inSink "github.com/mohammedimrankasab/metadata-ingestion-service/internal/sink"
	"go.uber.org/zap"
)

type Application struct {
	Components       *Components
	IngestionService *ingestion.Service
	Config           inConfig.Config
}

func NewApplication() (*Application, error) {

	log, err := inLog.NewLogger()
	if err != nil {
		return nil, err
	}
	metricsServer := server.New(log)
	config := inConfig.Config{
		WorkerCount:  4,
		JobQueueSize: 100,
	}

	powerBI := connectors.NewPowerBIConnector(log)
	consoleSink := inSink.NewConsoleSink(log)

	processor := processor.NewProcessor(
		log,
		consoleSink,
	)
	service := ingestion.New(
		log,
		config,
		processor,
		powerBI,
	)
	components := &Components{
		Logger:     log,
		Processor:  processor,
		Connectors: []connectors.Connector{powerBI},
		Sink:       consoleSink,
		Server:     metricsServer,
	}
	return &Application{
		Components:       components,
		IngestionService: service,
		Config:           config,
	}, nil
}

type Components struct {
	Logger     *zap.Logger
	Processor  *processor.Processor
	Connectors []connectors.Connector
	Sink       inSink.Sink
	Server     *server.Server
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
