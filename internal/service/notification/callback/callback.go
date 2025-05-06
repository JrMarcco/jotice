package callback

import (
	"github.com/JrMarcco/easy-kit/xsync"
	clientv1 "github.com/JrMarcco/jotice-api/api/client/v1"
	"github.com/JrMarcco/jotice/internal/domain"
	innergrpc "github.com/JrMarcco/jotice/internal/pkg/grpc"
	"github.com/JrMarcco/jotice/internal/pkg/logger"
	"github.com/JrMarcco/jotice/internal/repository"
	"github.com/JrMarcco/jotice/internal/service/config"
	"google.golang.org/grpc"
)

var _ Service = (*DefaultCallbackService)(nil)

type Service interface {
}

type DefaultCallbackService struct {
	configSvc     config.Service
	bizIdToConfig xsync.Map[uint64, *domain.CallbackConfig]
	clients       *innergrpc.Clients[clientv1.CallbackServiceClient]
	repo          repository.CallbackLogRepo
	logger        logger.Logger
}

func NewDefaultCallbackService(
	configSvc config.Service,
	repo repository.CallbackLogRepo,
	logger logger.Logger,
) *DefaultCallbackService {
	return &DefaultCallbackService{
		configSvc:     configSvc,
		bizIdToConfig: xsync.Map[uint64, *domain.CallbackConfig]{},
		repo:          repo,
		clients: innergrpc.NewClients(func(conn *grpc.ClientConn) clientv1.CallbackServiceClient {
			return clientv1.NewCallbackServiceClient(conn)
		}),
		logger: logger,
	}
}
