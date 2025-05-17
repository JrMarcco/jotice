package sender

import (
	"context"

	"github.com/JrMarcco/jotice/internal/domain"
	"github.com/JrMarcco/jotice/internal/pkg/logger"
	"github.com/JrMarcco/jotice/internal/repository"
	"github.com/JrMarcco/jotice/internal/service/channel"
	"github.com/JrMarcco/jotice/internal/service/config"
	"github.com/JrMarcco/jotice/internal/service/notification/callback"
)

//go:generate mockgen -source=./type.go -destination=./mock/sender.mock.go -package=sendermock -type=Sender
type Sender interface {
	Send(ctx context.Context, n domain.Notification) (domain.SendResp, error)
	BatchSend(ctx context.Context, ns []domain.Notification) ([]domain.SendResp, error)
}

type DefaultSender struct {
	repo         repository.NotificationRepo
	bizConfigSvc config.Service
	callbackSvc  callback.Service
	channel      channel.Channel
	logger       logger.Logger
}
