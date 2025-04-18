package config

import (
	"github.com/JrMarcco/jotice/internal/repository"
)

type BizConfigService interface{}

type DefaultBizConfigService struct {
	repo repository.BizConfigRepo
}

func NewDefaultBizConfigService(repo repository.BizConfigRepo) *DefaultBizConfigService {
	return &DefaultBizConfigService{
		repo: repo,
	}
}
