package connectors

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewCSVConnector(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	connector := NewCSVConnector(logger)

	if connector == nil {
		t.Fatal("expected connector to be created")
	}
}

func TestCSVConnector_Name(t *testing.T) {
	t.Parallel()

	connector := NewCSVConnector(zap.NewNop())

	if connector.Name() != "CSV" {
		t.Fatalf("expected CSV, got %q", connector.Name())
	}
}

func TestCSVConnector_FetchMetadata(t *testing.T) {
	t.Parallel()

	connector := NewCSVConnector(zap.NewNop())

	metadata, err := connector.FetchMetadata(
		context.Background(),
		nil,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metadata) != 1 {
		t.Fatalf("expected 1 metadata, got %d", len(metadata))
	}

	m := metadata[0]

	if m.Name != "employees.csv" {
		t.Errorf("expected employees.csv, got %s", m.Name)
	}

	if m.Type != "TABLE" {
		t.Errorf("expected TABLE, got %s", m.Type)
	}

	if m.Workspace != "Local" {
		t.Errorf("expected Local, got %s", m.Workspace)
	}

	if m.Source != "CSV" {
		t.Errorf("expected CSV, got %s", m.Source)
	}

	if m.ID == "" {
		t.Error("expected generated ID")
	}

	if m.LastModified.IsZero() {
		t.Error("expected LastModified to be set")
	}
}

func TestCSVConnector_FetchMetadata_LastSyncNil(t *testing.T) {
	t.Parallel()

	connector := NewCSVConnector(zap.NewNop())

	metadata, err := connector.FetchMetadata(
		context.Background(),
		nil,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metadata) != 1 {
		t.Fatalf("expected 1 metadata, got %d", len(metadata))
	}
}

func TestCSVConnector_FetchMetadata_LastSyncOld(t *testing.T) {
	t.Parallel()

	connector := NewCSVConnector(zap.NewNop())

	old := time.Now().Add(-1 * time.Hour)

	metadata, err := connector.FetchMetadata(
		context.Background(),
		&old,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metadata) != 1 {
		t.Fatalf("expected 1 metadata, got %d", len(metadata))
	}
}

func TestCSVConnector_FetchMetadata_LastSyncFuture(t *testing.T) {
	t.Parallel()

	connector := NewCSVConnector(zap.NewNop())

	future := time.Now().Add(1 * time.Hour)

	metadata, err := connector.FetchMetadata(
		context.Background(),
		&future,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metadata) != 0 {
		t.Fatalf("expected 0 metadata, got %d", len(metadata))
	}
}

func TestCSVConnector_Health(t *testing.T) {
	t.Parallel()

	connector := NewCSVConnector(zap.NewNop())

	err := connector.Health(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCSVConnector_Health_ContextCancelled(t *testing.T) {
	t.Parallel()

	connector := NewCSVConnector(zap.NewNop())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := connector.Health(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
