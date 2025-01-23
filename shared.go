package go_nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"log/slog"
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

type ServerOpts struct {
	instanceID         string
	timeout            time.Duration
	WithoutLeaderFns   bool
	WithoutFollowerFns bool
}

type ServerOption func(opts *ServerOpts, config *micro.Config)

func WithStatsHandler(handler micro.StatsHandler) ServerOption {
	return func(_ *ServerOpts, config *micro.Config) {
		config.StatsHandler = handler
	}
}

func WithDoneHandler(handler micro.DoneHandler) ServerOption {
	return func(_ *ServerOpts, config *micro.Config) {
		config.DoneHandler = handler
	}
}

func WithErrorHandler(handler micro.ErrHandler) ServerOption {
	return func(_ *ServerOpts, config *micro.Config) {
		config.ErrorHandler = handler
	}
}

func WithServiceVersion(version string) ServerOption {
	return func(_ *ServerOpts, config *micro.Config) {
		config.Version = version
	}
}

func WithoutLeaderFns() ServerOption {
	return func(opts *ServerOpts, _ *micro.Config) {
		opts.WithoutLeaderFns = true
	}
}

func WithoutFollowerFns() ServerOption {
	return func(opts *ServerOpts, _ *micro.Config) {
		opts.WithoutFollowerFns = true
	}
}

func NewService(name string, conn *nats.Conn, impl any, opts ...ServerOption) (micro.Service, *ServerOpts, error) {
	config := micro.Config{
		Name:    name,
		Version: "1.0.0",
	}

	// Check if the service implements any of the handler interfaces
	// but do it before applying options, so these can still override the handlers
	if statsHandler, isStatsHandler := impl.(StatsHandler); isStatsHandler {
		config.StatsHandler = statsHandler.Stats
		slog.Debug("Service implements StatsHandler; using service's Stats method")
	}
	if doneHandler, isDoneHandler := impl.(DoneHandler); isDoneHandler {
		config.DoneHandler = doneHandler.Done
		slog.Debug("Service implements DoneHandler; using service's Done method")
	}
	if errHandler, isErrHandler := impl.(ErrHandler); isErrHandler {
		config.ErrorHandler = errHandler.Err
		slog.Debug("Service implements ErrHandler; using service's Err method")
	}

	// Apply options
	options := new(ServerOpts)
	for _, opt := range opts {
		opt(options, &config)
	}

	// Create the service
	service, err := micro.AddService(conn, config)
	if err != nil {
		return nil, nil, err
	}

	return service, options, nil
}
