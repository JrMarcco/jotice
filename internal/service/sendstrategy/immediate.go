package sendstrategy

import (
	"context"
	"github.com/JrMarcco/jotice/internal/repository"

	"github.com/JrMarcco/jotice/internal/domain"
)

var _ SendStrategy = (*ImmediateSendStrategy)(nil)

type ImmediateSendStrategy struct {
	repo repository.NotificationRepo
}

func (s *ImmediateSendStrategy) Send(ctx context.Context, n domain.Notification) (domain.SendResp, error) {
	panic("not implemented")
}

func (s *ImmediateSendStrategy) BatchSend(ctx context.Context, ns []domain.Notification) ([]domain.SendResp, error) {
	panic("not implemented")
}
