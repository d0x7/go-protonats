package go_nats

import (
	"context"
	"time"
)

// CallOptions to  be used when calling a method.
type CallOptions interface {
	SetInstanceID(string)
	SetTimeout(time.Duration)
	SetRetry(int, time.Duration, time.Duration, context.Context)
}

type CallOption func(options CallOptions)

// WithInstanceID routes the call to the specified instance.
func WithInstanceID(id string) CallOption {
	return func(options CallOptions) {
		options.SetInstanceID(id)
	}
}

// WithTimeout overrides the default timeout for the call.
func WithTimeout(timeout time.Duration) CallOption {
	return func(options CallOptions) {
		options.SetTimeout(timeout)
	}
}

// WithRetry sets the number of retries, the minimum wait time, the maximum wait time, and the context for the call.
// Only used when NATS returns a NoResponders error, more or less efficiently "queueing" calls.
func WithRetry(ctx context.Context, minWait, maxWait time.Duration, maxTries int) CallOption {
	return func(opts CallOptions) {
		opts.SetRetry(maxTries, minWait, maxWait, ctx)
	}
}
