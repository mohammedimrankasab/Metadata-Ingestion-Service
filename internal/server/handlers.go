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
func (s *Server) Ingest(w http.ResponseWriter, r *http.Request) {

	go func() {
		if err := s.ingestion.Run(r.Context()); err != nil {
			s.logger.Error("ingestion failed", zap.Error(err))
		}
	}()

	w.WriteHeader(http.StatusAccepted)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Ingestion started",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
