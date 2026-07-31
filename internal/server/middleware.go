// Package server provides the HTTP server,
// middleware and REST API endpoints.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
		rw := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(rw, r)

		requestID, _ := r.Context().Value(RequestIDKey).(string)

		s.logger.Info(
			"http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
			zap.String("request_id", requestID),
			zap.Int("status", rw.status),
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

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
func (s *Server) Tracing(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			ctx, span := otel.Tracer(
				"metadata-ingestion-service/http",
			).Start(
				r.Context(),
				"HTTP "+r.Method+" "+r.URL.Path,
			)

			defer span.End()

			span.SetAttributes(
				attribute.String(
					"http.method",
					r.Method,
				),
				attribute.String(
					"http.route",
					r.URL.Path,
				),
			)

			r = r.WithContext(ctx)

			rw := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(
				rw,
				r,
			)

			span.SetAttributes(
				attribute.Int(
					"http.status_code",
					rw.status,
				),
			)

			if rw.status >= 500 {

				span.SetStatus(
					codes.Error,
					"HTTP request failed",
				)
			}
		},
	)
}
