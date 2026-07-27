package ingestion

import (
	"context"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/processor"
	"go.uber.org/zap"
)

type fakeConnector struct {
	metadata []models.Metadata
	err      error
	called   atomic.Bool
}

func (f *fakeConnector) Name() string {
	return "fake"
}

func (f *fakeConnector) Health(ctx context.Context) error {
	return nil
}

func (f *fakeConnector) FetchMetadata(
	ctx context.Context,
	lastSync *time.Time,
) ([]models.Metadata, error) {

	f.called.Store(true)

	return f.metadata, f.err
}

type fakeSink struct {
	writeCount atomic.Int32
}

func (f *fakeSink) Write(
	ctx context.Context,
	m models.Metadata,
) error {

	f.writeCount.Add(1)

	return nil
}
func TestRunProcessesAllMetadata(t *testing.T) {

	sink := &fakeSink{}

	p := processor.NewProcessor(
		zap.NewNop(),
		sink,
	)

	connector := &fakeConnector{
		metadata: []models.Metadata{
			models.NewMetadata(
				"1",
				"Dashboard",
				models.DashboardType,
				"Finance",
				"PowerBI",
				time.Now(),
			),
			models.NewMetadata(
				"2",
				"Report",
				models.ReportType,
				"Finance",
				"PowerBI",
				time.Now(),
			),
		},
	}

	cfg := &config.Config{
		WorkerCount:  2,
		JobQueueSize: 10,
	}

	service := New(
		zap.NewNop(),
		cfg,
		p,
		connector,
	)

	if err := service.Run(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !connector.called.Load() {
		t.Fatal("connector was not called")
	}

	if sink.writeCount.Load() != 2 {
		t.Fatalf(
			"expected 2 writes got %d",
			sink.writeCount.Load(),
		)
	}
}
func TestRunNoConnectors(t *testing.T) {
	sink := &fakeSink{}

	p := processor.NewProcessor(
		zap.NewNop(),
		sink,
	)

	cfg := &config.Config{
		WorkerCount:  1,
		JobQueueSize: 1,
	}

	service := New(
		zap.NewNop(),
		cfg,
		p,
	)

	err := service.Run(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "no connectors configured" {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestRunContextCancelled(t *testing.T) {
	sink := &fakeSink{}

	p := processor.NewProcessor(
		zap.NewNop(),
		sink,
	)

	var metadata []models.Metadata

	for i := 0; i < 10000; i++ {
		metadata = append(metadata,
			models.NewMetadata(
				strconv.Itoa(i),
				"Dashboard",
				models.DashboardType,
				"Finance",
				"PowerBI",
				time.Now(),
			),
		)
	}

	connector := &fakeConnector{
		metadata: metadata,
	}

	cfg := &config.Config{
		WorkerCount:  1,
		JobQueueSize: 1,
	}

	service := New(
		zap.NewNop(),
		cfg,
		p,
		connector,
	)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	err := service.Run(ctx)

	if err == nil {
		t.Fatal("expected context cancellation")
	}
}
