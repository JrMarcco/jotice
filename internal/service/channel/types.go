package channel

import (
	"context"

	"github.com/JrMarcco/jotice/internal/domain"
	"github.com/JrMarcco/jotice/internal/errs"
)

type Channel interface {
	Send(ctx context.Context, notification domain.Notification) (domain.SendResp, error)
}

var _ Channel = (*Dispatcher)(nil)

// Dispatcher is a channel dispatcher that chooses the appropriate channel based on the notification's channel configuration.
// Is a dispatcher pattern implementation.
type Dispatcher struct {
	channels map[domain.Channel]Channel
}

func (d *Dispatcher) Send(ctx context.Context, notification domain.Notification) (domain.SendResp, error) {
	ch, ok := d.channels[notification.Channel]
	if !ok {
		return domain.SendResp{}, errs.ErrInvalidChannel
	}

	return ch.Send(ctx, notification)
}

var _ Channel = (*baseChannel)(nil)

type baseChannel struct{}

func (b *baseChannel) Send(ctx context.Context, notification domain.Notification) (domain.SendResp, error) {
	//TODO implement me
	panic("implement me")
}
