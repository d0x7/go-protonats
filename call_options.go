package go_nats

import (
	"context"
	"time"
)

// CallOptions to be used when calling a method.
type CallOptions interface {
	SetInstanceID(string)
	SetTimeout(time.Duration)
	SetRetry(int, time.Duration, time.Duration, context.Context)
	WithoutFinisher()
	SetExtraSubject(subject string)
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

// WithoutFinisher disables the use of a finisher, which would return early
// after the first response when no further response is received after 250 ms
func WithoutFinisher() CallOption {
	return func(options CallOptions) {
		options.WithoutFinisher()
	}
}

// WithExtraSubject sets an extra subject for the subjects used by the microservice.
// Primarily used for consensus algorithms to distinguish between different services, while using the same implementation.
func WithExtraSubject(extraSubject string) CallOption {
	return func(options CallOptions) {
		options.SetExtraSubject(extraSubject)
	}
}
