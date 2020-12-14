package stack

import (
	"context"
	"time"

	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/cmd"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/transport"
)

type Options struct {
	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	Broker    broker.Options
	Cmd       cmd.Options
	Client    client.Options
	Server    server.Options
	Registry  registry.Options
	Transport transport.Options
	Selector  selector.Options
	Logger    logger.Options
	Config    config.Options

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	Signal bool
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Broker:    broker.Options{},
		Cmd:       cmd.Options{},
		Client:    client.Options{},
		Server:    server.Options{},
		Registry:  registry.Options{},
		Transport: transport.Options{},
		Selector:  selector.Options{},
		Logger:    logger.Options{},
		Config:    config.Options{},
		Context:   context.Background(),
		Signal:    true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Logger(l logger.Options) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func Broker(b broker.Options) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

func Cmd(c cmd.Options) Option {
	return func(o *Options) {
		o.Cmd = c
	}
}

func Client(c client.Options) Option {
	return func(o *Options) {
		o.Client = c
	}
}

func Config(c config.Options) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.Signal = b
	}
}

func Server(s server.Options) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// Registry sets the registry for the service
// and the underlying components
func Registry(r registry.Options) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Selector sets the selector for the service client
func Selector(s selector.Options) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

// Transport sets the transport for the service
// and the underlying components
func Transport(t transport.Options) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Convenience options

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) {
		o.Server.Address = addr
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Server.Name = n
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Server.Version = v
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Server.Metadata = md
	}
}

func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.Cmd.Flags = append(o.Cmd.Flags, flags...)
	}
}

func Action(a func(*cli.Context)) Option {
	return func(o *Options) {
		o.Cmd.Action = a
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.Server.RegisterTTL = t
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.Server.RegisterInterval = t
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		o.Client.Wrappers = w
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		o.Client.CallOptions.CallWrappers = w
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		o.Server.HdlrWrappers = w
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w ...server.SubscriberWrapper) Option {
	return func(o *Options) {
		o.Server.SubWrappers = w
	}
}

// Before and Afters

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
