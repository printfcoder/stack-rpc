package stack

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/cmd"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/debug/profile"
	"github.com/stack-labs/stack-rpc/debug/profile/pprof"
	"github.com/stack-labs/stack-rpc/debug/service/handler"
	"github.com/stack-labs/stack-rpc/env"
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/transport"
	"github.com/stack-labs/stack-rpc/util/wrapper"
)

type service struct {
	opts Options
	ctx  *stackContext

	cmd       cmd.Cmd
	broker    broker.Broker
	client    client.Client
	server    server.Server
	registry  registry.Registry
	transport transport.Transport
	selector  selector.Selector
	config    config.Config
	logger    log.Logger

	once sync.Once
}

func newService(opts ...Option) Service {
	options := newOptions(opts...)

	// service name
	serviceName := options.Server.Options().Name

	// wrap client to inject From-Service header on any calls
	options.Client = wrapper.FromService(serviceName, options.Client)

	return &service{
		opts: options,
		// inject the context
		ctx: stackCtx,
	}
}

func (s *service) Name() string {
	return s.ctx.runtime.server.Options().Name
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *service) Init(opts ...Option) error {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	s.once.Do(func() {
		// setup the plugins
		for _, p := range strings.Split(os.Getenv(env.StackPlugin), ",") {
			if len(p) == 0 {
				continue
			}

			// load the plugin
			c, err := plugin.Load(p)
			if err != nil {
				log.Fatal(err)
			}

			// initialise the plugin
			if err := plugin.Init(c); err != nil {
				log.Fatal(err)
			}
		}
	})

	return nil
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Client() client.Client {
	return s.ctx.runtime.client
}

func (s *service) Server() server.Server {
	return s.ctx.runtime.server
}

func (s *service) String() string {
	return "stack"
}

func (s *service) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Stop() error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := s.server.Stop(); err != nil {
		return err
	}

	if err := s.config.Close(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *service) Run() error {
	if err := s.cmd.Init(
		cmd.Broker(&s.broker),
		cmd.Registry(&s.registry),
		cmd.Transport(&s.transport),
		cmd.Client(&s.client),
		cmd.Server(&s.server),
		cmd.Selector(&s.selector),
		cmd.Logger(&s.logger),
		cmd.Config(&s.config),
	); err != nil {
		log.Errorf("cmd init error: %s", err)
		return err
	}

	// register the debug handler
	if err := s.server.Handle(
		s.server.NewHandler(
			handler.DefaultHandler,
			server.InternalHandler(true),
		),
	); err != nil {
		return err
	}

	// start the profiler
	// TODO: set as an option to the service, don't just use pprof
	if prof := os.Getenv(env.StackDebugProfile); len(prof) > 0 {
		service := s.server.Options().Name
		version := s.server.Options().Version
		id := s.server.Options().Id
		profiler := pprof.NewProfile(
			profile.Name(service + "." + version + "." + id),
		)
		if err := profiler.Start(); err != nil {
			return err
		}
		defer profiler.Stop()
	}

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	}

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}
