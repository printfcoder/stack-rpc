package stack

import (
	"context"
	"time"

	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/cmd"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/debug/profile"
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/transport"
)

type Options struct {
	Cmd         cmd.Cmd
	serviceOpts []service.Option
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Cmd: cmd.NewCmd(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.serviceOpts = append(o.serviceOpts, service.Lo
	}
}

func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Service.Broker = b
	}
}

func Cmd(c cmd.Cmd) Option {
	return func(o *Options) {
		o.Cmd = c
	}
}

func Client(c client.Client) Option {
	return func(o *Options) {
		o.Service.Client = c
	}
}

func Config(c config.Config) Option {
	return func(o *Options) {
		o.Service.Config = c
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Service.Context = ctx
	}
}

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.Service.Signal = b
	}
}

func Server(s server.Server) Option {
	return func(o *Options) {
		o.Service.Server = s
	}
}

// Registry sets the registry for the service
// and the underlying components
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Service.Registry = r
	}
}

// Selector sets the selector for the service client
func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Service.Selector = s
	}
}

// Transport sets the transport for the service
// and the underlying components
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Service.Transport = t
	}
}

// Convenience options

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) {
		o.Service.Server.Init(server.Address(addr))
	}
}

// Unique server id
func Id(id string) Option {
	return func(o *Options) {
		o.Service.Id = id
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Service.Name = n
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Service.Server.Init(server.Version(v))
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Service.Server.Init(server.Metadata(md))
	}
}

func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.Cmd.App().Flags = append(o.Cmd.App().Flags, flags...)
	}
}

func Action(a func(*cli.Context)) Option {
	return func(o *Options) {
		o.Cmd.App().Action = a
	}
}

// Profile to be used for debug profile
func Profile(p profile.Profile) Option {
	return func(o *Options) {
		o.Service.Profile = p
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.Service.Server.Init(server.RegisterTTL(t))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.Service.Server.Init(server.RegisterInterval(t))
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			o.Service.Client = w[i-1](o.Client)
		}
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		o.Service.Client.Init(client.WrapCall(w...))
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// Init once
		o.Service.Server.Init(wrappers...)
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w ...server.SubscriberWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}

		// Init once
		o.Service.Server.Init(wrappers...)
	}
}

// Before and Afters

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.Service.BeforeStart = append(o.Service.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.Service.BeforeStop = append(o.Service.BeforeStop, fn)
	}
}

func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.Service.AfterStart = append(o.Service.AfterStart, fn)
	}
}

func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.Service.AfterStop = append(o.Service.AfterStop, fn)
	}
}
