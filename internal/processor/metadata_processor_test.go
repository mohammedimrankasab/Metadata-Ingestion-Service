package processor

import (
	"context"
	"testing"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type MockSink struct {
	Called bool
}

func (m *MockSink) Write(
	ctx context.Context,
	metadata models.Metadata,
) error {

	m.Called = true

	return nil
}
func TestProcessorProcess(t *testing.T) {

	logger := zap.NewNop()

	mockSink := &MockSink{}

	p := NewProcessor(
		logger,
		mockSink,
	)

	job := models.NewJob(
		"PowerBI",
		models.NewMetadata(
			"id",
			"name",
			models.ReportType,
			"default",
			"PowerBI",
			time.Now(),
		),
	)

	err := p.Process(context.Background(), job)

	require.NoError(t, err)

	require.True(t, mockSink.Called)
}
