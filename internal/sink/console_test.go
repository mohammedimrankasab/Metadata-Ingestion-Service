package sink

import (
	"context"
	"testing"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

func TestNewConsoleSink(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	sink := NewConsoleSink(logger)

	if sink == nil {
		t.Fatal("expected sink to be created")
	}

	if sink.logger != logger {
		t.Fatal("logger was not assigned")
	}
}

func TestConsoleSink_Write(t *testing.T) {
	t.Parallel()

	sink := NewConsoleSink(zap.NewNop())

	metadata := models.NewMetadata(
		"1",
		"Sales Dashboard",
		models.DashboardType,
		"Finance",
		"PowerBI",
		time.Now(),
	)

	err := sink.Write(
		context.Background(),
		metadata,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestConsoleSink_Write_CancelledContext(t *testing.T) {
	t.Parallel()

	sink := NewConsoleSink(zap.NewNop())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	metadata := models.NewMetadata(
		"1",
		"Sales Dashboard",
		models.DashboardType,
		"Finance",
		"PowerBI",
		time.Now(),
	)

	err := sink.Write(ctx, metadata)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
