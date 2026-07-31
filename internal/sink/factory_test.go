package sink

import (
	"testing"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"go.uber.org/zap"
)

func TestNew_ConsoleSink(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		SinkType: "console",
	}

	sink, err := New(cfg, zap.NewNop())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sink == nil {
		t.Fatal("expected sink")
	}

	if _, ok := sink.(*ConsoleSink); !ok {
		t.Fatalf("expected *ConsoleSink, got %T", sink)
	}
}

func TestNew_OpenSearchSink(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		SinkType:        "opensearch",
		OpenSearchURL:   "http://localhost:9200",
		OpenSearchIndex: "metadata",
	}

	sink, err := New(cfg, zap.NewNop())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sink == nil {
		t.Fatal("expected sink")
	}

	osSink, ok := sink.(*OpenSearchSink)
	if !ok {
		t.Fatalf("expected *OpenSearchSink, got %T", sink)
	}

	if osSink.url != cfg.OpenSearchURL {
		t.Fatal("url mismatch")
	}

	if osSink.index != cfg.OpenSearchIndex {
		t.Fatal("index mismatch")
	}
}

func TestNew_JSONSink_NotImplemented(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		SinkType: "json",
	}

	sink, err := New(cfg, zap.NewNop())

	if err == nil {
		t.Fatal("expected error")
	}

	if sink != nil {
		t.Fatal("expected nil sink")
	}
}

func TestNew_InvalidSink(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		SinkType: "abc",
	}

	sink, err := New(cfg, zap.NewNop())

	if err == nil {
		t.Fatal("expected error")
	}

	if sink != nil {
		t.Fatal("expected nil sink")
	}
}
