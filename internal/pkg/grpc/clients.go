package grpc

import (
	"fmt"

	"github.com/JrMarcco/easy-kit/xsync"
	"google.golang.org/grpc"
)

type Clients[T any] struct {
	clientMap xsync.Map[string, T]
	creator   func(conn *grpc.ClientConn) T
}

func (cs *Clients[T]) Get(serviceName string) (T, error) {
	if client, ok := cs.clientMap.Load(serviceName); ok {
		return client, nil
	}

	conn, err := grpc.NewClient(fmt.Sprintf("etcd:///%s", serviceName))
	if err != nil {
		var zero T
		return zero, fmt.Errorf("create client failed: %w", err)
	}

	client := cs.creator(conn)
	cs.clientMap.Store(serviceName, client)
	return client, nil
}

func NewClients[T any](creator func(conn *grpc.ClientConn) T) *Clients[T] {
	return &Clients[T]{creator: creator}
}
