package repository

import (
	"github.com/JrMarcco/jotice/internal/pkg/logger"
	"github.com/JrMarcco/jotice/internal/repository/cache"
	"github.com/JrMarcco/jotice/internal/repository/dao"
)

type BizConfigRepo interface{}

type DefaultBizConfigRepo struct {
	dao    dao.ConfigDAO
	lc     cache.BizConfigCache
	rc     cache.BizConfigCache
	logger logger.Logger
}
