package local

import (
	"github.com/JrMarcco/jotice/internal/pkg/logger"
	"github.com/JrMarcco/jotice/internal/repository/cache"
	gcache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

var _ cache.BizConfigCache = (*LBizConfigCache)(nil)

// LBizConfigCache is a local cache implementation for biz config.
type LBizConfigCache struct {
	c      *gcache.Cache
	rdb    redis.Cmdable
	logger logger.Logger
}
