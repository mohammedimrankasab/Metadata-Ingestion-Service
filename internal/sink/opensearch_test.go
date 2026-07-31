package sink

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/models"
	"go.uber.org/zap"
)

func TestNewOpenSearchSink(t *testing.T) {
	t.Parallel()

	sink := NewOpenSearchSink(
		zap.NewNop(),
		"http://localhost:9200",
		"metadata",
	)

	if sink == nil {
		t.Fatal("expected sink")
	}

	if sink.client == nil {
		t.Fatal("expected http client")
	}

	if sink.url != "http://localhost:9200" {
		t.Fatal("incorrect url")
	}

	if sink.index != "metadata" {
		t.Fatal("incorrect index")
	}
}

func TestOpenSearchSink_Write_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			t.Fatalf("expected POST got %s", r.Method)
		}

		if r.URL.Path != "/metadata/_doc" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("unexpected content type %s", ct)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "Sales Dashboard") {
			t.Fatal("metadata not marshalled")
		}

		w.WriteHeader(http.StatusCreated)

	}))
	defer server.Close()

	sink := NewOpenSearchSink(
		zap.NewNop(),
		server.URL,
		"metadata",
	)

	metadata := models.NewMetadata(
		"id",
		"Sales Dashboard",
		models.DashboardType,
		"Finance",
		"PowerBI",
		time.Now(),
	)

	err := sink.Write(context.Background(), metadata)

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestOpenSearchSink_Write_ServerError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusInternalServerError)

	}))
	defer server.Close()

	sink := NewOpenSearchSink(
		zap.NewNop(),
		server.URL,
		"metadata",
	)

	metadata := models.NewMetadata(
		"id",
		"Sales Dashboard",
		models.DashboardType,
		"Finance",
		"PowerBI",
		time.Now(),
	)

	err := sink.Write(context.Background(), metadata)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOpenSearchSink_Write_InvalidURL(t *testing.T) {
	t.Parallel()

	sink := NewOpenSearchSink(
		zap.NewNop(),
		"://invalid",
		"metadata",
	)

	metadata := models.NewMetadata(
		"id",
		"Sales Dashboard",
		models.DashboardType,
		"Finance",
		"PowerBI",
		time.Now(),
	)

	err := sink.Write(context.Background(), metadata)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOpenSearchSink_Write_ContextCancelled(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		time.Sleep(200 * time.Millisecond)

	}))
	defer server.Close()

	sink := NewOpenSearchSink(
		zap.NewNop(),
		server.URL,
		"metadata",
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	metadata := models.NewMetadata(
		"id",
		"Sales Dashboard",
		models.DashboardType,
		"Finance",
		"PowerBI",
		time.Now(),
	)

	err := sink.Write(ctx, metadata)

	if err == nil {
		t.Fatal("expected error")
	}
}
