package config

import (
	"os"
	"runtime"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("WORKER_COUNT", "")
	t.Setenv("JOB_QUEUE_SIZE", "")
	t.Setenv("HTTP_PORT", "")

	cfg := Load()

	if cfg.WorkerCount != runtime.NumCPU() {
		t.Fatalf("expected %d got %d", runtime.NumCPU(), cfg.WorkerCount)
	}

	if cfg.JobQueueSize != 100 {
		t.Fatalf("unexpected queue size")
	}

	if cfg.HTTPPort != "8080" {
		t.Fatalf("unexpected http port")
	}
}

func TestLoadEnvironment(t *testing.T) {
	t.Setenv("WORKER_COUNT", "5")
	t.Setenv("JOB_QUEUE_SIZE", "50")
	t.Setenv("HTTP_PORT", "9000")

	cfg := Load()

	if cfg.WorkerCount != 5 {
		t.Fail()
	}

	if cfg.JobQueueSize != 50 {
		t.Fail()
	}

	if cfg.HTTPPort != "9000" {
		t.Fail()
	}
}

func TestInvalidWorkerCountFallsBack(t *testing.T) {

	os.Setenv("WORKER_COUNT", "-5")

	cfg := Load()

	if cfg.WorkerCount != runtime.NumCPU() {
		t.Fail()
	}
}
