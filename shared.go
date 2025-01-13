package go_nats

import (
	"github.com/nats-io/nats.go/micro"
	"time"
)

type Ping struct {
	micro.ServiceIdentity
	Type string
	RTT  time.Duration
}

type CallOpts struct {
	instanceID string
	timeout    time.Duration
}

func (opts *CallOpts) HasInstanceID() bool {
	return opts.instanceID != ""
}

func (opts *CallOpts) HasTimeout() bool {
	return opts.timeout != 0
}

func (opts *CallOpts) GetInstanceID() string {
	return opts.instanceID
}

func (opts *CallOpts) GetTimeout() time.Duration {
	return opts.timeout
}

func (opts *CallOpts) GetTimeoutOr(duration time.Duration) time.Duration {
	if opts.HasTimeout() {
		return opts.GetTimeout()
	}
	return duration
}

type CallOption func(*CallOpts)

func WithInstanceID(id string) CallOption {
	return func(opts *CallOpts) {
		opts.instanceID = id
	}
}

func WithTimeout(timeout time.Duration) CallOption {
	return func(options *CallOpts) {
		options.timeout = timeout
	}
}

func ProcessCallOptions(opts ...CallOption) CallOpts {
	options := CallOpts{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type ServerOption func(config *micro.Config)

func WithStatsHandler(handler micro.StatsHandler) ServerOption {
	return func(config *micro.Config) {
		config.StatsHandler = handler
	}
}

func WithDoneHandler(handler micro.DoneHandler) ServerOption {
	return func(config *micro.Config) {
		config.DoneHandler = handler
	}
}

func WithErrorHandler(handler micro.ErrHandler) ServerOption {
	return func(config *micro.Config) {
		config.ErrorHandler = handler
	}
}
