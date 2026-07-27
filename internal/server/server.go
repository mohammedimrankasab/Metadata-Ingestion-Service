// Package server provides the HTTP server,
// middleware and REST API endpoints.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/ingestion"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type IngestionService interface {
	Run(ctx context.Context) error
}

type Server struct {
	logger    *zap.Logger
	server    *http.Server
	cfg       *config.Config
	ingestion IngestionService
}

func New(logger *zap.Logger, cfg *config.Config, ingestion *ingestion.Service) *Server {

	s := &Server{
		logger:    logger,
		ingestion: ingestion,
		cfg:       cfg,
	}

	mux := http.NewServeMux()

	mux.Handle("/health", s.withMiddleware(http.HandlerFunc(s.Health)))
	mux.Handle("/ready", s.withMiddleware(http.HandlerFunc(s.Ready)))
	mux.Handle("/ingest", s.withMiddleware(http.HandlerFunc(s.Ingest)))
	mux.Handle("/metrics", promhttp.Handler())

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return s
}

func (s *Server) Start() {

	go func() {

		s.logger.Info(
			"HTTP server started",
			zap.String("addr", s.server.Addr),
		)

		if err := s.server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			s.logger.Error(
				"HTTP server stopped",
				zap.Error(err),
			)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {

	s.logger.Info("Shutting down HTTP server")

	return s.server.Shutdown(ctx)
}
func (s *Server) withMiddleware(h http.Handler) http.Handler {
	return s.Recovery(
		s.Logging(
			RequestID(h),
		),
	)
}
