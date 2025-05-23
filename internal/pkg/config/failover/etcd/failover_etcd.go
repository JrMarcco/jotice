package etcd

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/JrMarcco/jotice/internal/errs"
	"github.com/JrMarcco/jotice/internal/pkg/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const failoverPrefix = "/config/failover"

var _ config.FailoverManager = (*ManagerOfEtcd)(nil)

type ManagerOfEtcd struct {
	mu sync.RWMutex

	client      *clientv3.Client
	watchCancel []context.CancelFunc
}

func (m *ManagerOfEtcd) Failover(ctx context.Context, si config.ServiceInstance) error {
	if err := si.Validate(); err != nil {
		return err
	}

	_, err := m.client.Put(ctx, m.serviceKey(si), "")
	return err
}

func (m *ManagerOfEtcd) Recover(ctx context.Context, si config.ServiceInstance) error {
	if err := si.Validate(); err != nil {
		return err
	}

	_, err := m.client.Delete(ctx, m.serviceKey(si))
	return err
}

func (m *ManagerOfEtcd) WatchFailover(ctx context.Context) (<-chan config.FailoverEvent, error) {
	watchCtx, cancel := context.WithCancel(ctx)
	m.mu.Lock()
	m.watchCancel = append(m.watchCancel, cancel)
	m.mu.Unlock()

	watchChan := m.client.Watch(watchCtx, failoverPrefix, clientv3.WithPrefix())
	eventChan := make(chan config.FailoverEvent)

	go func() {
		for {
			select {
			case resp := <-watchChan:
				if resp.Canceled {
					return
				}

				if err := resp.Err(); err != nil {
					continue
				}

				for _, event := range resp.Events {
					if event.Type == clientv3.EventTypePut && string(event.Kv.Value) != "" {
						key := string(event.Kv.Key)
						si, err := m.parseKey(key)
						if err != nil {
							continue
						}

						copySi := config.ServiceInstance{
							Name:  si.Name,
							Group: si.Group,
							Addr:  si.Addr,
						}
						select {
						case eventChan <- config.FailoverEvent{Si: copySi}:
						case <-watchCtx.Done():
							return
						}
					}
				}
			case <-watchCtx.Done():
				return
			}
		}
	}()

	return eventChan, nil
}

func (m *ManagerOfEtcd) parseKey(key string) (*config.ServiceInstance, error) {
	splits := strings.Split(strings.Trim(strings.TrimPrefix(key, failoverPrefix), "/"), "/")

	if len(splits) != 3 {
		return nil, fmt.Errorf("%w: invalid etcd service instance key: %s", errs.ErrInvalidParam, key)
	}

	return &config.ServiceInstance{
		Name:  splits[0],
		Group: splits[1],
		Addr:  splits[2],
	}, nil
}

func (m *ManagerOfEtcd) WatchRecover(ctx context.Context, si config.ServiceInstance) (<-chan struct{}, error) {
	if err := si.Validate(); err != nil {
		return nil, err
	}

	key := m.serviceKey(si)

	watchCtx, cancel := context.WithCancel(ctx)
	m.mu.Lock()
	m.watchCancel = append(m.watchCancel, cancel)
	m.mu.Unlock()

	watchChan := m.client.Watch(watchCtx, key)
	notifyChan := make(chan struct{})

	go func() {
		defer close(notifyChan)

		for {
			select {
			case resp := <-watchChan:
				if resp.Canceled {
					return
				}

				if err := resp.Err(); err != nil {
					continue
				}

				for _, event := range resp.Events {
					if event.Type == clientv3.EventTypeDelete {
						select {
						case notifyChan <- struct{}{}:
						case <-watchCtx.Done():
							return
						}
					}
				}
			case <-watchCtx.Done():
				return
			}
		}
	}()

	return notifyChan, nil
}

func (m *ManagerOfEtcd) TryTakeover(ctx context.Context, undertakenSi, targetSi config.ServiceInstance) (bool, error) {
	if err := undertakenSi.Validate(); err != nil {
		return false, err
	}

	if err := targetSi.Validate(); err != nil {
		return false, err
	}

	siKey := m.serviceKey(undertakenSi)
	resp, err := m.client.Get(ctx, siKey)
	if err != nil {
		return false, err
	}

	if len(resp.Kvs) == 0 {
		// no service needs to be take over
		return false, fmt.Errorf("%w", errs.ErrNoAvailableFailoverService)
	}

	if string(resp.Kvs[0].Value) != "" {
		return false, fmt.Errorf("%w: service %s is already undertaken", errs.ErrNoAvailableFailoverService, undertakenSi.Name)
	}

	txnResp, err := m.client.Txn(ctx).
		If(clientv3.Compare(clientv3.Value(siKey), "=", "")).
		Then(clientv3.OpPut(siKey, m.serviceVal(undertakenSi))).
		Commit()
	if err != nil {
		return false, err
	}

	return txnResp.Succeeded, nil
}

func (m *ManagerOfEtcd) serviceKey(si config.ServiceInstance) string {
	return fmt.Sprintf("%s/%s/%s%s", failoverPrefix, si.Name, si.Group, si.Addr)
}

func (m *ManagerOfEtcd) serviceVal(si config.ServiceInstance) string {
	return fmt.Sprintf("%s:%s:%s", si.Name, si.Group, si.Addr)
}

func (m *ManagerOfEtcd) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, cancel := range m.watchCancel {
		cancel()
	}
	m.watchCancel = nil
	return nil
}

func NewManager(client *clientv3.Client) *ManagerOfEtcd {
	return &ManagerOfEtcd{
		client:      client,
		watchCancel: make([]context.CancelFunc, 0),
	}
}
