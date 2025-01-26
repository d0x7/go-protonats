package go_nats

import (
	"context"
	"time"
)

type CallOptions interface {
	SetInstanceID(string)
	SetTimeout(time.Duration)
	SetRetry(int, time.Duration, time.Duration, context.Context)
}

type CallOption func(options CallOptions)

func WithInstanceID(id string) CallOption {
	return func(options CallOptions) {
		options.SetInstanceID(id)
	}
}

func WithTimeout(timeout time.Duration) CallOption {
	return func(options CallOptions) {
		options.SetTimeout(timeout)
	}
}

func WithRetry(ctx context.Context, minWait, maxWait time.Duration, maxTries int) CallOption {
	return func(opts CallOptions) {
		opts.SetRetry(maxTries, minWait, maxWait, ctx)
	}
}
