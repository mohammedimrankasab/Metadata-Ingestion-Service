// Package server provides the HTTP server,
// middleware and REST API endpoints.
package server

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "UP",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) Ready(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "READY",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) Ingest(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	s.ingestMu.Lock()

	if s.running {

		s.ingestMu.Unlock()

		w.WriteHeader(
			http.StatusConflict,
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"message": "Ingestion already running",
			},
		)

		return
	}

	s.running = true

	s.ingestMu.Unlock()

	go func() {

		defer func() {

			s.ingestMu.Lock()

			s.running = false

			s.ingestMu.Unlock()

		}()

		if err := s.ingestion.Run(
			s.appCtx,
		); err != nil {

			s.logger.Error(
				"ingestion failed",
				zap.Error(err),
			)
		}

	}()

	w.WriteHeader(
		http.StatusAccepted,
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "Ingestion started",
		},
	)
}
