package impl

import (
	"context"
	"time"
	"xiam.li/go-nats"
)

type CallOpts struct {
	InstanceID      string
	Timeout         time.Duration
	Retries         int
	RetryDelay      time.Duration
	RetryContext    context.Context
	DisableFinisher bool
	ExtraSubject    string
}

func (opts *CallOpts) WithoutFinisher() {
	opts.DisableFinisher = true
}

func ProcessCallOptions(opts ...go_nats.CallOption) *CallOpts {
	options := new(CallOpts)
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func (opts *CallOpts) SetInstanceID(id string) {
	opts.InstanceID = id
}

func (opts *CallOpts) SetTimeout(timeout time.Duration) {
	opts.Timeout = timeout
}

func (opts *CallOpts) SetRetry(maxTries int, minWait time.Duration, maxWait time.Duration, ctx context.Context) {
	if maxTries <= 0 {
		return
	}

	retries := float64(maxTries)
	minWaitf := float64(minWait.Milliseconds())
	maxWaitf := float64(maxWait.Milliseconds())

	totalMinWait := retries * minWaitf
	if totalMinWait > maxWaitf {
		panic("retries multiplied by minWait must be below maxWait")
	}

	extraTime := maxWaitf - totalMinWait
	timePerRetry := minWaitf + (extraTime / retries)

	opts.Retries = maxTries
	opts.RetryContext = ctx
	opts.RetryDelay = time.Duration(timePerRetry) * time.Millisecond
}

func (opts *CallOpts) HasInstanceID() bool {
	return opts.InstanceID != ""
}

func (opts *CallOpts) HasTimeout() bool {
	return opts.Timeout != 0
}

func (opts *CallOpts) GetTimeoutOr(duration time.Duration) time.Duration {
	if opts.HasTimeout() {
		return opts.Timeout
	}
	return duration
}

func (opts *CallOpts) ShouldRetry() bool {
	return opts.Retries > 0 && opts.RetryContext != nil
}

func (opts *CallOpts) SetExtraSubject(subject string) {
	opts.ExtraSubject = subject
}

func (opts *CallOpts) Subject(subject string) string {
	return _subject(subject, opts.ExtraSubject, opts.InstanceID)
}

// Interface guard
var _ go_nats.CallOptions = (*CallOpts)(nil)
