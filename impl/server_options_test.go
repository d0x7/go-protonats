package impl

import (
	"github.com/nats-io/nats.go/micro"
	"reflect"
	"testing"
	go_nats "xiam.li/go-protonats"
)

func TestServerOpts_WithServiceVersion(t *testing.T) {
	opts := new(ServerOpts)
	go_nats.WithServiceVersion("test")(opts)
	if opts.ServiceVersion != "test" {
		t.Error("ServiceVersion not set correctly")
	}
}

func TestServerOpts_WithoutFollowerFns(t *testing.T) {
	opts := new(ServerOpts)
	go_nats.WithoutFollowerFns()(opts)
	if !opts.WithoutFollowerFunctions {
		t.Error("WithoutFollowerFunctions not set correctly")
	}
}

func TestServerOpts_WithoutLeaderFns(t *testing.T) {
	opts := new(ServerOpts)
	go_nats.WithoutLeaderFns()(opts)
	if !opts.WithoutLeaderFunctions {
		t.Error("WithoutLeaderFunctions not set correctly")
	}
}

func TestServerOpts_WithExtraSubject(t *testing.T) {
	opts := new(ServerOpts)
	go_nats.WithExtraSubjectSrv("test")(opts)
	if opts.ExtraSubject != "test" {
		t.Error("ExtraSubject not set correctly")
	}
}

func TestServerOpts_Subject(t *testing.T) {
	const empty = ""
	t.Run("WithoutExtraOrSuffix", func(t *testing.T) {
		if _subject("topic", empty, empty) != "topic" {
			t.Error("Subject not set correctly")
		}
	})
	t.Run("WithExtra", func(t *testing.T) {
		if _subject("topic", "extra", empty) != "topic.extra" {
			t.Error("Subject not set correctly")
		}
	})
	t.Run("WithSuffix", func(t *testing.T) {
		if _subject("topic", empty, "suffix") != "topic.suffix" {
			t.Error("Subject not set correctly")
		}
	})
	t.Run("WithExtraAndSuffix", func(t *testing.T) {
		if _subject("topic", "extra", "suffix") != "topic.extra.suffix" {
			t.Error("Subject not set correctly")
		}
	})
}

func TestProcessServerOptions(t *testing.T) {
	t.Run("WithMicroConfig", func(t *testing.T) {
		statsHandler := func(endpoint *micro.Endpoint) any {
			return "stats_test"
		}
		doneHandler := func(_ micro.Service) {}
		errHandler := func(_ micro.Service, _ *micro.NATSError) {}
		config := new(micro.Config)
		opts := ProcessServerOptions(
			config,
			go_nats.WithServiceVersion("test"),
			go_nats.WithStatsHandler(statsHandler),
			go_nats.WithDoneHandler(doneHandler),
			go_nats.WithErrorHandler(errHandler),
		)
		if config.Version != "test" {
			t.Error("ServiceVersion not set correctly")
		}
		if config.StatsHandler == nil {
			t.Error("StatsHandler not set correctly")
		}
		statsHandlerPtr := reflect.ValueOf(statsHandler).Pointer()
		configStatsHandlerPtr := reflect.ValueOf(config.StatsHandler).Pointer()
		if statsHandlerPtr != configStatsHandlerPtr {
			t.Error("StatsHandler not set correctly")
		}
		if (*opts.StatsHandler)(nil) != "stats_test" {
			t.Error("StatsHandler not set correctly")
		}
		if config.DoneHandler == nil {
			t.Error("DoneHandler not set correctly")
		}
		doneHandlerPtr := reflect.ValueOf(doneHandler).Pointer()
		configDoneHandlerPtr := reflect.ValueOf(config.DoneHandler).Pointer()
		if doneHandlerPtr != configDoneHandlerPtr {
			t.Error("DoneHandler not set correctly")
		}
		if config.ErrorHandler == nil {
			t.Error("ErrorHandler not set correctly")
		}
		errHandlerPtr := reflect.ValueOf(errHandler).Pointer()
		configErrHandlerPtr := reflect.ValueOf(config.ErrorHandler).Pointer()
		if errHandlerPtr != configErrHandlerPtr {
			t.Error("ErrorHandler not set correctly")
		}
	})

	t.Run("WithServerOpts", func(t *testing.T) {

	})
}
