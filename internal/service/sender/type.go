package sender

import (
	"context"
	"github.com/JrMarcco/jotice/internal/domain"
	"github.com/JrMarcco/jotice/internal/repository"
	"github.com/JrMarcco/jotice/internal/service/channel"
)

type Sender interface {
	Send(ctx context.Context, n domain.Notification) (domain.SendResp, error)
	BatchSend(ctx context.Context, ns []domain.Notification) ([]domain.SendResp, error)
}

type DefaultSender struct {
	repo         repository.NotificationRepo
	bizConfigSvc repository.BizConfigRepo
	channel      channel.Channel
}
