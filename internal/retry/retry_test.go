package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryEventuallySucceeds(t *testing.T) {

	attempts := 0

	cfg := Config{
		MaxRetries: 3,
		BaseDelay:  time.Millisecond,
	}

	err := Do(context.Background(), cfg, func() error {

		attempts++

		if attempts < 3 {
			return errors.New("failed")
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts got %d", attempts)
	}
}

func TestRetryFails(t *testing.T) {

	cfg := Config{
		MaxRetries: 2,
		BaseDelay:  time.Millisecond,
	}

	err := Do(context.Background(), cfg, func() error {
		return errors.New("boom")
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRetryCancelled(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := Config{
		MaxRetries: 5,
		BaseDelay:  time.Millisecond,
	}

	err := Do(ctx, cfg, func() error {
		return nil
	})

	if err == nil {
		t.Fatal("expected context error")
	}
}
