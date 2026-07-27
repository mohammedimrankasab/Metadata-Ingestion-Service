package connectors

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestPowerBIFetchMetadata(t *testing.T) {

	c := NewPowerBIConnector(zap.NewNop())

	data, err := c.FetchMetadata(context.Background(), nil)

	if err != nil {
		t.Fatal(err)
	}

	if len(data) != sampleMetadataCount {
		t.Fatalf("expected %d got %d", sampleMetadataCount, len(data))
	}
}

func TestPowerBIHealth(t *testing.T) {

	c := NewPowerBIConnector(zap.NewNop())

	if err := c.Health(context.Background()); err != nil {
		t.Fatal(err)
	}
}
