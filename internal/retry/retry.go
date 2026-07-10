package retry

import (
	"context"
	"errors"
	"time"
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

var ErrPermanent = errors.New("permanent error")
