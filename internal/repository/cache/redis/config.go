package redis

import (
	"github.com/JrMarcco/jotice/internal/pkg/logger"
	"github.com/JrMarcco/jotice/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

var _ cache.ConfigCache = (*RCacheConfig)(nil)

type RCacheConfig struct {
	rdb    *redis.Cmdable
	logger logger.Logger
}
