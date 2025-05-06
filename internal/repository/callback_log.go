package repository

import (
	"context"
	"github.com/JrMarcco/jotice/internal/domain"
	"github.com/JrMarcco/jotice/internal/repository/dao"
)

var _ CallbackLogRepo = (*DefaultCallbackLogRepo)(nil)

type CallbackLogRepo interface {
	ListByNotificationIds(ctx context.Context, ids []uint64) ([]domain.CallbackLog, error)
}

type DefaultCallbackLogRepo struct {
	dao dao.CallbackLogDAO
}

func (d *DefaultCallbackLogRepo) ListByNotificationIds(ctx context.Context, ids []uint64) ([]domain.CallbackLog, error) {
	panic("implement me")
}
