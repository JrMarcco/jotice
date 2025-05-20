package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/JrMarcco/easy-kit/xsync"
	"github.com/JrMarcco/jotice/internal/pkg/sharding"
	"github.com/JrMarcco/jotice/internal/pkg/snowflake"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Notification entity definition.
type Notification struct {
	Id            uint64
	BizId         uint64
	BizKey        string
	Receivers     string
	Channel       string
	TplId         uint64
	TplVersionId  uint64
	TplParams     string
	Status        string
	ScheduleStrat int64
	ScheduleEnd   int64
	Version       int32
	CreatedAt     int64
	UpdatedAt     int64
}

type NotificationDAO interface {
	GetByBizKey(ctx context.Context, bizId uint64, bizKey string) (Notification, error)
	GetByBizKeys(ctx context.Context, bizId uint64, bizKeys ...string) ([]Notification, error)

	FindDreadyNotifications(ctx context.Context, offset, limit int) ([]Notification, error)
}

var _ NotificationDAO = (*NotifShardingDAO)(nil)

type NotifShardingDAO struct {
	dbs *xsync.Map[string, *gorm.DB]

	notifShardingStrategy sharding.Strategy
	cbLogShardingStrategy sharding.Strategy

	idGenerator *snowflake.Generator
}

func (n *NotifShardingDAO) GetByBizKey(ctx context.Context, bizId uint64, bizKey string) (Notification, error) {
	dst := n.notifShardingStrategy.Shard(bizId, bizKey)
	dstDB, ok := n.dbs.Load(dst.DB)
	if !ok {
		return Notification{}, fmt.Errorf("unknown db: %s", dst.DB)
	}

	var notif Notification
	err := dstDB.WithContext(ctx).Table(dst.Table).
		Where("`biz_id` = ? AND `biz_key` = ?", bizId, bizKey).
		First(&notif).Error
	if err != nil {
		return Notification{}, fmt.Errorf("failed to get notification by BizId = %d and BizKey = %s, cause of: %w", bizId, bizKey, err)
	}
	return notif, nil
}

func (n *NotifShardingDAO) GetByBizKeys(ctx context.Context, bizId uint64, bizKeys ...string) ([]Notification, error) {
	//notifMap := make(map[[2]string][]string, len(bizKeys))
	//for index := range bizKeys {
	//	bizKey := bizKeys[index]
	//
	//	dst := n.notifShardingStrategy.Shard(bizId, bizKey)
	//
	//	shardingInfo := [2]string{dst.DB, dst.Table}
	//
	//	val, ok := notifMap[shardingInfo]
	//	if !ok {
	//		val = []string{bizKey}
	//	} else {
	//		val = append(val, bizKey)
	//	}
	//	notifMap[shardingInfo] = val
	//}
	//
	//var eg errgroup.Group
	//for shardingInfo, ks := range notifMap {
	//	eg.Go(func() error {
	//		dbName := shardingInfo[0]
	//		tableName := shardingInfo[1]
	//
	//		gormDB, ok := n.dbs.Load(dbName)
	//		if !ok {
	//			return fmt.Errorf("unknown db: %s", dbName)
	//		}
	//
	//		var notifs []Notification
	//		err := gormDB.WithContext(ctx).Table(tableName).
	//			Where("`biz_id` = ? AND `biz_key` IN (?)", bizId, ks).
	//			Find(&notifs).Error
	//		if err != nil {
	//			return err
	//		}
	//	})
	//}

	//TODO implement me
	panic("implement me")
}

func (n *NotifShardingDAO) FindDreadyNotifications(ctx context.Context, offset, limit int) ([]Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (n *NotifShardingDAO) isUniqueConstraintErr(err error) bool {
	if err == nil {
		return false
	}

	mysqlErr := new(mysql.MySQLError)
	if ok := errors.As(err, &mysqlErr); ok {
		const uniqueConstraintErrCode = 1062
		return mysqlErr.Number == uniqueConstraintErrCode
	}
	return false
}

func NewNotifShardingDAO(
	dbs *xsync.Map[string, *gorm.DB],
	notifShardingStrategy sharding.Strategy,
	cbLogShardingStrategy sharding.Strategy,
	idGenerator *snowflake.Generator,
) *NotifShardingDAO {
	return &NotifShardingDAO{
		dbs:                   dbs,
		notifShardingStrategy: notifShardingStrategy,
		cbLogShardingStrategy: cbLogShardingStrategy,
		idGenerator:           idGenerator,
	}
}
