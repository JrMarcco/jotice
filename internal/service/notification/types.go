package notification

import (
	"context"

	"github.com/JrMarcco/jotice/internal/domain"
)

//go:generate mockgen -source=./type.go -destination=./mock/service.mock.go -package=notificationmock -type=Service
type Service interface{}

//go:generate mockgen -source=./type.go -destination=./mock/send_service.mock.go -package=notificationmock -type=SendService
type SendService interface {
	Send(ctx context.Context, n domain.Notification) (domain.SendResp, error)
	AsyncSend(ctx context.Context, n domain.Notification) (domain.SendResp, error)
	BatchSend(ctx context.Context, ns []domain.Notification) (domain.BatchSendResp, error)
	BatchAsyncSend(ctx context.Context, ns []domain.Notification) (domain.BatchAsyncSendResp, error)
}
