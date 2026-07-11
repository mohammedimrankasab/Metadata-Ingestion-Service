package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/config"
	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/ingestion"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	logger    *zap.Logger
	server    *http.Server
	cfg       config.Config
	ingestion *ingestion.Service
}

func New(logger *zap.Logger, cfg config.Config, ingestion *ingestion.Service) *Server {

	s := &Server{
		logger:    logger,
		ingestion: ingestion,
		cfg:       cfg,
	}

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	handler := s.Recovery(
		s.Logging(
			RequestID(
				http.HandlerFunc(s.Health),
			),
		),
	)

	mux.Handle("/health", handler)

	mux.HandleFunc("/ready", s.Ready)
	mux.HandleFunc("/ingest", s.Ingest)

	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", s.cfg.MetricsPort),
			Handler: mux,
		},
		ingestion: ingestion,
	}
}

func (s *Server) Start() {

	go func() {

		s.logger.Info(
			"Metrics server started",
			zap.String("addr", s.server.Addr),
		)

		if err := s.server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			s.logger.Error(
				"Metrics server stopped",
				zap.Error(err),
			)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {

	s.logger.Info("Shutting down metrics server")

	return s.server.Shutdown(ctx)
}
