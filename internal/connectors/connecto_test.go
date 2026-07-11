package connectors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPowerBIFetchMetadata(t *testing.T) {

	logger := zap.NewNop()

	connector := NewPowerBIConnector(logger)

	result, err := connector.FetchMetadata(
		context.Background(),
		nil,
	)

	require.NoError(t, err)

	require.Len(t, result, 200)
}
