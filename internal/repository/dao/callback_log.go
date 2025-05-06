package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type CallbackLog struct {
	Id             uint64 `gorm:"column:id"`
	NotificationId uint64 `gorm:"column:notification_id"`
	RetryTimes     int32  `gorm:"column:retry_times"`
	NextRetryAt    int64  `gorm:"column:next_retry_at"`
	Status         string `gorm:"column:status"`
	CreatedAt      int64  `gorm:"column:created_at"`
	UpdatedAt      int64  `gorm:"column:updated_at"`
}

func (c CallbackLog) TableName() string {
	return "callback_log"
}

var _ CallbackLogDAO = (*DefaultCallbackLogDAO)(nil)

type CallbackLogDAO interface {
	BatchUpdate(ctx context.Context, logs []CallbackLog) error
	ListPendingBatch(ctx context.Context, startTime int64, startId uint64, batchSize int32) ([]CallbackLog, uint64, error)
	ListByNotificationIds(ctx context.Context, ids []uint64) ([]CallbackLog, error)
}

type DefaultCallbackLogDAO struct {
	db *gorm.DB
}

func (d *DefaultCallbackLogDAO) BatchUpdate(ctx context.Context, logs []CallbackLog) error {
	if len(logs) == 0 {
		return nil
	}

	updateAt := time.Now().UnixMilli()
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, log := range logs {
			res := tx.Model(&CallbackLog{Id: log.Id}).
				Updates(map[string]any{
					"retry_times":   log.RetryTimes,
					"next_retry_at": log.NextRetryAt,
					"status":        log.Status,
					"updated_at":    updateAt,
				})

			if res.Error != nil {
				return res.Error
			}
		}
		return nil
	})
}

func (d *DefaultCallbackLogDAO) ListPendingBatch(
	ctx context.Context, startTime int64, startId uint64, batchSize int32,
) ([]CallbackLog, uint64, error) {
	var logs []CallbackLog

	nextStartId := 0
	res := d.db.WithContext(ctx).Model(&CallbackLog{}).
		Where("status = ?", "pending").
		Where("next_retry_at <= ?", startTime).
		Where("id > ?", startId).
		Order("create_at ASC").
		Limit(int(batchSize)).
		Find(&logs)

	if res.Error != nil {
		return nil, uint64(nextStartId), res.Error
	}

	if len(logs) > 0 {
		nextStartId = int(logs[len(logs)-1].Id)
	}

	return logs, uint64(nextStartId), nil
}

func (d *DefaultCallbackLogDAO) ListByNotificationIds(ctx context.Context, ids []uint64) ([]CallbackLog, error) {
	var logs []CallbackLog
	err := d.db.WithContext(ctx).Where("notification_id IN (?)", ids).Find(&logs).Error
	return logs, err
}

func NewDefaultCallbackLogDAO(db *gorm.DB) *DefaultCallbackLogDAO {
	return &DefaultCallbackLogDAO{
		db: db,
	}
}
