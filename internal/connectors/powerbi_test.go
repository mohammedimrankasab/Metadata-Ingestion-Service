package connectors

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewPowerBIConnector(t *testing.T) {
	t.Parallel()

	connector := NewPowerBIConnector(zap.NewNop())

	if connector == nil {
		t.Fatal("expected connector")
	}
}

func TestPowerBI_Name(t *testing.T) {
	t.Parallel()

	connector := NewPowerBIConnector(zap.NewNop())

	if connector.Name() != "PowerBI" {
		t.Fatalf("expected PowerBI got %s", connector.Name())
	}
}

func TestPowerBI_FetchMetadata(t *testing.T) {
	t.Parallel()

	connector := NewPowerBIConnector(zap.NewNop())

	data, err := connector.FetchMetadata(
		context.Background(),
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(data) != sampleMetadataCount {
		t.Fatalf(
			"expected %d got %d",
			sampleMetadataCount,
			len(data),
		)
	}
}

func TestPowerBI_FetchMetadata_LastSyncOld(t *testing.T) {
	t.Parallel()

	connector := NewPowerBIConnector(zap.NewNop())

	lastSync := time.Now().Add(-1 * time.Hour)

	data, err := connector.FetchMetadata(
		context.Background(),
		&lastSync,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(data) != sampleMetadataCount {
		t.Fatalf("expected %d got %d", sampleMetadataCount, len(data))
	}
}

func TestPowerBI_FetchMetadata_LastSyncFuture(t *testing.T) {
	t.Parallel()

	connector := NewPowerBIConnector(zap.NewNop())

	lastSync := time.Now().Add(1 * time.Hour)

	data, err := connector.FetchMetadata(
		context.Background(),
		&lastSync,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(data) != 0 {
		t.Fatalf("expected 0 got %d", len(data))
	}
}

func TestPowerBI_FetchMetadata_ContextCancelled(t *testing.T) {
	t.Parallel()

	connector := NewPowerBIConnector(zap.NewNop())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := connector.FetchMetadata(
		ctx,
		nil,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPowerBI_Health(t *testing.T) {
	t.Parallel()

	connector := NewPowerBIConnector(zap.NewNop())

	if err := connector.Health(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestPowerBI_Health_ContextCancelled(t *testing.T) {
	t.Parallel()

	connector := NewPowerBIConnector(zap.NewNop())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := connector.Health(ctx); err != nil {
		t.Fatal(err)
	}
}
