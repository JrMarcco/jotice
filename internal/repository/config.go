package repository

import (
	"github.com/JrMarcco/jotice/internal/pkg/logger"
	"github.com/JrMarcco/jotice/internal/repository/dao"
)

type BizConfigRepo interface{}

type DefaultBizConfigRepo struct {
	dao    dao.BizConfigDAO
	logger logger.Logger
}
