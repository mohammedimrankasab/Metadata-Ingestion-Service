package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestRequestID(t *testing.T) {

	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Context().Value(RequestIDKey) == nil {
			t.Fail()
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
}

func TestRecovery(t *testing.T) {

	s := &Server{
		logger: zap.NewNop(),
	}

	h := s.Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fail()
	}
}
