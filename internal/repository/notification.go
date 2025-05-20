package repository

import (
	"context"

	"github.com/JrMarcco/jotice/internal/domain"
	"github.com/JrMarcco/jotice/internal/repository/dao"
	"go.uber.org/zap"
)

// NotificationRepo is a repository for notification.
type NotificationRepo interface {
	GetByBizKey(ctx context.Context, bizId uint64, bizKey string) (domain.Notification, error)
	GetByBizKeys(ctx context.Context, bizId uint64, bizKeys ...string) ([]domain.Notification, error)

	FindDreadyNotifications(ctx context.Context, offset, limit int) ([]domain.Notification, error)
}

var _ NotificationRepo = (*DefaultNotifRepo)(nil)

// DefaultNotifRepo is a default implementation of NotificationRepo.
type DefaultNotifRepo struct {
	dao    dao.NotificationDAO
	logger *zap.Logger
}

func (d *DefaultNotifRepo) GetByBizKey(ctx context.Context, bizId uint64, bizKey string) (domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultNotifRepo) GetByBizKeys(ctx context.Context, bizId uint64, bizKeys ...string) ([]domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultNotifRepo) FindDreadyNotifications(ctx context.Context, offset, limit int) ([]domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func NewNotificationRepo(dao dao.NotificationDAO, logger *zap.Logger) *DefaultNotifRepo {
	return &DefaultNotifRepo{
		dao:    dao,
		logger: logger,
	}
}
