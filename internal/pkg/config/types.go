package config

import (
	"context"
	"fmt"

	"github.com/JrMarcco/jotice/internal/errs"
)

type ServiceInstance struct {
	Name  string
	Group string
	Addr  string
}

func (si *ServiceInstance) Validate() error {
	if si.Name == "" {
		return fmt.Errorf("%w: service instance name should not be empty", errs.ErrInvalidParam)
	}

	if si.Addr == "" {
		return fmt.Errorf("%W: invalidate service instance address: %s", errs.ErrInvalidParam, si.Addr)
	}
	return nil
}

type FailoverEvent struct {
	Si ServiceInstance
}

type FailoverManager interface {
	Failover(ctx context.Context, si ServiceInstance) error
	Recover(ctx context.Context, si ServiceInstance) error

	WatchFailover(ctx context.Context) (<-chan FailoverEvent, error)
	WatchRecover(ctx context.Context, si ServiceInstance) (<-chan struct{}, error)

	TryTakeover(ctx context.Context, undertakenSi, targetSi ServiceInstance) (bool, error)

	Close() error
}
