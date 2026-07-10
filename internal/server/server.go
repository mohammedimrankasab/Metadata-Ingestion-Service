package server

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	server *http.Server
}

func New(logger *zap.Logger) *Server {

	mux := http.NewServeMux()

	mux.Handle(
		"/metrics",
		promhttp.Handler(),
	)

	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:    ":2112",
			Handler: mux,
		},
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
