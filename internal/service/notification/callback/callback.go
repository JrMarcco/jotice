package callback

import (
	"github.com/JrMarcco/easy-kit/xsync"
	"github.com/JrMarcco/jotice/internal/domain"
	"github.com/JrMarcco/jotice/internal/service/config"
)

var _ Service = (*DefaultCallbackService)(nil)

type Service interface {
}

type DefaultCallbackService struct {
	configSvc     config.Service
	bizIdToConfig xsync.Map[uint64, *domain.CallbackConfig]
}
