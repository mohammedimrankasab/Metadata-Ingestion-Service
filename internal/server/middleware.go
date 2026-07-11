package server

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

func RequestID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := uuid.NewString()

		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(
			r.Context(),
			RequestIDKey,
			requestID,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}
func (s *Server) Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		next.ServeHTTP(w, r)

		s.logger.Info(
			"HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
func (s *Server) Recovery(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {

				s.logger.Error(
					"Panic recovered",
					zap.Any("panic", err),
				)

				http.Error(
					w,
					"Internal Server Error",
					http.StatusInternalServerError,
				)
			}

		}()

		next.ServeHTTP(w, r)
	})
}
