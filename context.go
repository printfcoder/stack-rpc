package stack

import (
	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/transport"
)

var stackCtx = &stackContext{}

type stackContext struct {
	runtime *runtime
}

type runtime struct {
	broker    broker.Broker
	registry  registry.Registry
	transport transport.Transport
	client    client.Client
	server    server.Server
	selector  selector.Selector
}

func (r *runtime) Registry() registry.Registry {
	return r.registry
}

func (r *runtime) Broker() broker.Broker {
	return r.broker
}

func (r *runtime) Transport() transport.Transport {
	return r.transport
}

func (r *runtime) Client() transport.Transport {
	return r.transport
}

func (r *runtime) Server() server.Server {
	return r.server
}

func (r *runtime) Selector() selector.Selector {
	return r.selector
}
