package impl

import (
	"context"
	"testing"
	"time"
	"xiam.li/go-nats"
)

func TestCallOpts_WithInstanceID(t *testing.T) {
	opts := new(CallOpts)
	go_nats.WithInstanceID("test")(opts)
	if opts.InstanceID != "test" {
		t.Error("InstanceID not set correctly")
	}
}

func TestCallOpts_WithTimeout(t *testing.T) {
	opts := new(CallOpts)
	go_nats.WithTimeout(100 * time.Millisecond)(opts)
	if opts.Timeout != 100*time.Millisecond {
		t.Error("Timeout not set correctly")
	}
}

func TestCallOpts_WithRetry(t *testing.T) {
	opts := new(CallOpts)
	go_nats.WithRetry(context.Background(), 100*time.Millisecond, 300*time.Millisecond, 3)(opts)
	if opts.Retries != 3 {
		t.Error("Retries not set correctly")
	}
	if opts.RetryDelay != 100*time.Millisecond {
		t.Error("RetryDelay not set correctly", opts.RetryDelay)
	}
	if opts.RetryContext != context.Background() {
		t.Error("RetryContext not set correctly")
	}
	t.Run("invalidMaxTries", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Invalid maxTries should panic")
			}
		}()
		opts := new(CallOpts)
		go_nats.WithRetry(context.Background(), 100*time.Millisecond, 300*time.Millisecond, 4)(opts)
		t.Error("Should have panicked")
	})
	t.Run("noRetries", func(t *testing.T) {
		opts := new(CallOpts)
		go_nats.WithRetry(context.Background(), 100*time.Millisecond, 300*time.Millisecond, 0)(opts)
		if opts.Retries != 0 {
			t.Error("Retries should be 0")
		}
		if opts.RetryDelay != 0 {
			t.Error("RetryDelay should be 0")
		}
		if opts.RetryContext != nil {
			t.Error("RetryContext should be nil")
		}
	})
}

func TestProcessCallOptions(t *testing.T) {
	opts := ProcessCallOptions(
		go_nats.WithInstanceID("test"),
		go_nats.WithTimeout(100*time.Millisecond),
		go_nats.WithRetry(context.Background(), 100*time.Millisecond, 300*time.Millisecond, 3),
	)
	if opts.InstanceID != "test" {
		t.Error("InstanceID not set correctly")
	}
	if opts.Timeout != 100*time.Millisecond {
		t.Error("Timeout not set correctly")
	}
	if opts.Retries != 3 {
		t.Error("Retries not set correctly")
	}
	if opts.RetryDelay != 100*time.Millisecond {
		t.Error("RetryDelay not set correctly", opts.RetryDelay)
	}
	if opts.RetryContext != context.Background() {
		t.Error("RetryContext not set correctly")
	}
}
