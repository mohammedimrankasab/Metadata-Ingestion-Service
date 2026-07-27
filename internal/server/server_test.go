package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"
)

type fakeIngestion struct {
	started chan struct{}
}

func (f *fakeIngestion) Run(ctx context.Context) error {
	close(f.started)
	return nil
}
func TestIngest(t *testing.T) {

	ingestion := &fakeIngestion{
		started: make(chan struct{}),
	}

	s := &Server{
		logger:    zap.NewNop(),
		ingestion: ingestion,
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/ingest",
		nil,
	)

	rec := httptest.NewRecorder()

	s.Ingest(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected %d got %d",
			http.StatusAccepted,
			rec.Code,
		)
	}

	select {
	case <-ingestion.started:
		// success

	case <-time.After(time.Second):
		t.Fatal("ingestion was not started")
	}
}
