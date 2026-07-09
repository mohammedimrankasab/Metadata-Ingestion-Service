package app

import (
	"context"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/connectors"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/ingestion"
	inLog "github.com/mohammedimrankasab/metadata-ingestion-service/internal/logger"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/processor"
	"go.uber.org/zap"
)

type Application struct {
	Components       *Components
	IngestionService *ingestion.Service
}

func NewApplication() (*Application, error) {

	log, err := inLog.NewLogger()
	if err != nil {
		return nil, err
	}

	powerBI := connectors.NewPowerBIConnector(log)

	service := ingestion.New(
		log,
		powerBI,
	)
	components := &Components{
		Logger:     log,
		Connectors: []connectors.Connector{powerBI},
	}
	return &Application{
		Components:       components,
		IngestionService: service,
	}, nil
}

type Components struct {
	Logger     *zap.Logger
	Processor  *processor.Processor
	Connectors []connectors.Connector
}

func (app *Application) Run(ctx context.Context) error {
	defer func() {
		_ = app.Components.Logger.Sync()
	}()
	return app.IngestionService.Run(ctx)
}
