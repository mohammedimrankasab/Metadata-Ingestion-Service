package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestInitTracer(t *testing.T) {
	t.Parallel()

	shutdown, err := InitTracer(false)
	if err != nil {
		t.Fatalf("InitTracer() returned error: %v", err)
	}

	if shutdown == nil {
		t.Fatal("expected shutdown function")
	}

	if otel.GetTracerProvider() == nil {
		t.Fatal("expected tracer provider to be initialized")
	}

	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown returned error: %v", err)
	}
}

func TestInitTracer_MultipleCalls(t *testing.T) {
	t.Parallel()

	shutdown1, err := InitTracer(false)
	if err != nil {
		t.Fatalf("InitTracer() returned error: %v", err)
	}

	shutdown2, err := InitTracer(false)
	if err != nil {
		t.Fatalf("InitTracer() returned error: %v", err)
	}

	if shutdown1 == nil || shutdown2 == nil {
		t.Fatal("expected shutdown functions")
	}

	if err := shutdown1(context.Background()); err != nil {
		t.Fatalf("shutdown1 returned error: %v", err)
	}

	if err := shutdown2(context.Background()); err != nil {
		t.Fatalf("shutdown2 returned error: %v", err)
	}
}

func TestServiceName(t *testing.T) {
	t.Parallel()

	const expected = "metadata-ingestion-service"

	if ServiceName != expected {
		t.Fatalf("expected %q, got %q", expected, ServiceName)
	}
}
