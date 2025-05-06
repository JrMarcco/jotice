package redis

import (
	"github.com/JrMarcco/jotice/internal/pkg/logger"
	"github.com/JrMarcco/jotice/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

var _ cache.BizConfigCache = (*RBizCacheConfig)(nil)

// RBizCacheConfig is a redis cache implementation for biz config.
type RBizCacheConfig struct {
	rdb    *redis.Cmdable
	logger logger.Logger
}
