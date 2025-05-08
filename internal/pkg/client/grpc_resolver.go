package client

import (
	"context"
	"time"

	"github.com/JrMarcco/jotice/internal/pkg/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

var _ resolver.Builder = (*grpcResolverBuilder)(nil)

type grpcResolverBuilder struct {
	registry registry.Registry
	timeout  time.Duration
}

func (b *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	gr := &grpcResolver{
		target:   target,
		subConn:  cc,
		registry: b.registry,
		timeout:  b.timeout,
		close:    make(chan struct{}, 1),
	}

	gr.resolve()
	go gr.watch()

	return gr, nil
}

func (b *grpcResolverBuilder) Scheme() string {
	return "registry"
}

func NewGrpcResolverBuilder(registry registry.Registry, timeout time.Duration) resolver.Builder {
	return &grpcResolverBuilder{
		registry: registry,
		timeout:  timeout,
	}
}

var _ resolver.Resolver = (*grpcResolver)(nil)

type grpcResolver struct {
	target   resolver.Target
	subConn  resolver.ClientConn
	registry registry.Registry
	timeout  time.Duration

	close chan struct{}
}

func (g *grpcResolver) ResolveNow(_ resolver.ResolveNowOptions) {
	// re get all service
	g.resolve()
}

func (g *grpcResolver) resolve() {
	serviceName := g.target.Endpoint()
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	instances, err := g.registry.ListService(ctx, serviceName)
	cancel()

	if err != nil {
		g.subConn.ReportError(err)
	}

	addrs := make([]resolver.Address, len(instances))
	for _, inst := range instances {
		addrs = append(addrs, resolver.Address{
			Addr:       inst.Address,
			ServerName: inst.Name,
			Attributes: attributes.New(attrNameReadWeight, inst.ReadWeight).
				WithValue(attrNameWriteWeight, inst.WriteWeight).
				WithValue(attrNameGroup, inst.Group).
				WithValue(attrNameNode, inst.Name),
		})
	}

	err = g.subConn.UpdateState(resolver.State{
		Addresses: addrs,
	})
	if err != nil {
		g.subConn.ReportError(err)
	}
}

func (g *grpcResolver) Close() {
	g.close <- struct{}{}
}

func (g *grpcResolver) watch() {
	events := g.registry.Subscribe(g.target.Endpoint())

	for {
		select {
		case <-events:
			g.resolve()
		case <-g.close:
			return
		}
	}
}
