package notification

import (
	"context"
	"fmt"

	"github.com/JrMarcco/jotice/internal/domain"
	"github.com/JrMarcco/jotice/internal/errs"
	"github.com/JrMarcco/jotice/internal/repository"
)

//go:generate mockgen -source=./types.go -destination=./mock/service.mock.go -package=notificationmock -type=Service
type Service interface {
	// FindReadyNotifications find notifications that are ready to be schedule to send.
	FindReadyNotifications(ctx context.Context, offset, limit int) ([]domain.Notification, error)
	// GetByBizKeys get notifications by biz id and biz keys.
	GetByBizKeys(ctx context.Context, BizId uint64, bizKeys ...string) ([]domain.Notification, error)
}

var _ Service = (*DefaultNotifService)(nil)

type DefaultNotifService struct {
	repo repository.NotificationRepo
}

func (d *DefaultNotifService) FindReadyNotifications(ctx context.Context, offset, limit int) ([]domain.Notification, error) {
	return d.repo.FindDreadyNotifications(ctx, offset, limit)
}

func (d *DefaultNotifService) GetByBizKeys(ctx context.Context, BizId uint64, bizKeys ...string) ([]domain.Notification, error) {
	if len(bizKeys) == 0 {
		return nil, fmt.Errorf("%w: business keys should not be empty", errs.ErrInvalidParam)
	}

	notifications, err := d.repo.GetByBizKeys(ctx, BizId, bizKeys...)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications by biz keys, cause of: %w", err)
	}
	return notifications, nil
}

func NewDefaultNotifService(repo repository.NotificationRepo) *DefaultNotifService {
	return &DefaultNotifService{
		repo: repo,
	}
}
