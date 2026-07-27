// Package retry provides a reusable retry mechanism with
// exponential backoff and context cancellation support.
package retry

import (
	"context"
	"time"

	"github.com/mohammedimrankasab/metadata-ingestion-service/internal/metrics"
)

type Config struct {
	MaxRetries int
	BaseDelay  time.Duration
}

func Do(
	ctx context.Context,
	cfg Config,
	fn func() error,
) error {

	var err error

	delay := cfg.BaseDelay

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = fn()

		if err == nil {
			return nil
		}

		if attempt == cfg.MaxRetries {
			break
		}

		// Record that another retry will be attempted.
		metrics.RetryCount.Inc()

		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()

		case <-timer.C:
		}

		delay *= 2
	}

	return err
}
