package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetrySuccessFirstAttempt(t *testing.T) {

	cfg := Config{
		MaxRetries: 3,
		BaseDelay:  10 * time.Millisecond,
	}

	calls := 0

	err := Do(context.Background(), cfg, func() error {
		calls++
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 1, calls)
}

func TestRetrySucceedsAfterRetry(t *testing.T) {

	cfg := Config{
		MaxRetries: 3,
		BaseDelay:  5 * time.Millisecond,
	}

	attempts := 0

	err := Do(context.Background(), cfg, func() error {

		attempts++

		if attempts < 3 {
			return errors.New("temporary")
		}

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 3, attempts)
}

func TestRetryFails(t *testing.T) {

	cfg := Config{
		MaxRetries: 3,
		BaseDelay:  5 * time.Millisecond,
	}

	attempts := 0

	err := Do(context.Background(), cfg, func() error {

		attempts++

		return errors.New("failed")

	})

	require.Error(t, err)
	require.Equal(t, 4, attempts)
}
func TestRetryCancelled(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	cfg := Config{
		MaxRetries: 3,
		BaseDelay:  time.Second,
	}

	err := Do(ctx, cfg, func() error {
		return errors.New("boom")
	})

	require.Error(t, err)
}
