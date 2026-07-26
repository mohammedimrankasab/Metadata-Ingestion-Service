// Package logger provides application-wide structured logging
// using Uber's Zap logger.
package logger

import (
	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	return zap.NewProduction()
}
