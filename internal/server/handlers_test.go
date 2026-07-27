package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestHealth(t *testing.T) {

	s := &Server{
		logger: zap.NewNop(),
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	rec := httptest.NewRecorder()

	s.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fail()
	}
}

func TestReady(t *testing.T) {

	s := &Server{
		logger: zap.NewNop(),
	}

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)

	rec := httptest.NewRecorder()

	s.Ready(rec, req)

	if rec.Code != http.StatusOK {
		t.Fail()
	}
}
