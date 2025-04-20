package callback

import (
	"github.com/JrMarcco/jotice/internal/service/config"
)

var _ Service = (*DefaultCallbackService)(nil)

type Service interface {
}

type DefaultCallbackService struct {
	configSvc config.BizConfigService
	//bizIdToConfig sync.Map[int64, *domain.CallbackConfig]
}
