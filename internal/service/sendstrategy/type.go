package sendstrategy

import (
	"context"
	"fmt"

	"github.com/JrMarcco/jotice/internal/domain"
	"github.com/JrMarcco/jotice/internal/errs"
)

type SendStrategy interface {
	// Send send notification use strategy in the notification's strategy configuration.
	Send(ctx context.Context, n domain.Notification) (domain.SendResp, error)
	// BatchSend batch sends notifications.
	BatchSend(ctx context.Context, ns []domain.Notification) ([]domain.SendResp, error)
}

var _ SendStrategy = (*Dispatcher)(nil)

// Dispatcher is a strategy dispatcher that chooses the appropriate strategy based on the notification's strategy configuration.
type Dispatcher struct {
	defaultStrategy   *DefaultSendStrategy
	immediateStrategy *ImmediateSendStrategy
}

func (d *Dispatcher) Send(ctx context.Context, n domain.Notification) (domain.SendResp, error) {
	return d.chooseStrategy(n).Send(ctx, n)
}

func (d *Dispatcher) BatchSend(ctx context.Context, ns []domain.Notification) ([]domain.SendResp, error) {
	if len(ns) == 0 {
		return nil, fmt.Errorf("%w: no notifications to send", errs.ErrInvalidParam)
	}

	return d.chooseStrategy(ns[0]).BatchSend(ctx, ns)
}

func (d *Dispatcher) chooseStrategy(n domain.Notification) SendStrategy {
	if n.StrategyConfig.Strategy == domain.SendStrategyImmediate {
		return d.immediateStrategy
	}

	return d.defaultStrategy
}

func NewDispatcher(defaultStrategy *DefaultSendStrategy, immediateStrategy *ImmediateSendStrategy) *Dispatcher {
	return &Dispatcher{
		defaultStrategy:   defaultStrategy,
		immediateStrategy: immediateStrategy,
	}
}
