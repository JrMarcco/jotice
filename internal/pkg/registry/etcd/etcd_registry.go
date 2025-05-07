package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/JrMarcco/jotice/internal/pkg/registry"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var typeMap = map[mvccpb.Event_EventType]registry.EventType{
	mvccpb.PUT:    registry.EventTypePut,
	mvccpb.DELETE: registry.EventTypeDelete,
}

var _ registry.Registry = (*RegistryOfEtcd)(nil)

type RegistryOfEtcd struct {
	mu sync.RWMutex

	session     *concurrency.Session
	client      *clientv3.Client
	watchCancel []func()
}

func (r *RegistryOfEtcd) Register(ctx context.Context, si registry.ServiceInstance) error {
	val, err := json.Marshal(si)
	if err != nil {
		return err
	}

	_, err = r.client.Put(
		ctx, r.instanceKey(si), string(val), clientv3.WithLease(r.session.Lease()),
	)
	return err
}

func (r *RegistryOfEtcd) UnRegister(ctx context.Context, si registry.ServiceInstance) error {
	_, err := r.client.Delete(ctx, r.instanceKey(si))
	return err
}

func (r *RegistryOfEtcd) instanceKey(si registry.ServiceInstance) string {
	return fmt.Sprintf("/notification/%s/%s", si.Name, si.Address)
}

func (r *RegistryOfEtcd) ListService(ctx context.Context, serviceName string) ([]registry.ServiceInstance, error) {
	resp, err := r.client.Get(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	res := make([]registry.ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var si registry.ServiceInstance
		if err := json.Unmarshal(kv.Value, &si); err != nil {
			return nil, err
		}

		res = append(res, si)
	}

	return res, nil
}

func (r *RegistryOfEtcd) serviceKey(serviceName string) string {
	return fmt.Sprintf("/notification/%s", serviceName)
}

func (r *RegistryOfEtcd) Subscribe(serviceName string) <-chan registry.Event {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = clientv3.WithRequireLeader(ctx)

	r.mu.Lock()
	r.watchCancel = append(r.watchCancel, cancel)
	r.mu.Unlock()

	ch := r.client.Watch(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
	res := make(chan registry.Event)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case resp := <-ch:
				if resp.Canceled {
					return
				}

				if resp.Err() != nil {
					continue
				}

				for _, e := range resp.Events {
					res <- registry.Event{
						Type: typeMap[e.Type],
					}
				}
			}
		}
	}()

	return res
}

func (r *RegistryOfEtcd) Close() error {
	r.mu.Lock()
	for _, cancel := range r.watchCancel {
		cancel()
	}
	r.mu.Unlock()

	// can't close r.client.
	// because it will be used by other goroutines
	return r.session.Close()
}

func NewRegistry(client *clientv3.Client) (*RegistryOfEtcd, error) {
	session, err := concurrency.NewSession(client)
	if err != nil {
		return nil, err
	}

	return &RegistryOfEtcd{
		session: session,
		client:  client,
	}, nil
}
