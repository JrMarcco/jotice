package notification

import (
	"context"

	"github.com/JrMarcco/jotice/internal/domain"
)

type SendService interface {
	Send(ctx context.Context, n domain.Notification) (domain.SendResp, error)
	AsyncSend(ctx context.Context, n domain.Notification) (domain.SendResp, error)
	BatchSend(ctx context.Context, ns []domain.Notification) (domain.BatchSendResp, error)
	BatchAsyncSend(ctx context.Context, ns []domain.Notification) (domain.BatchAsyncSendResp, error)
}
