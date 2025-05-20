package notification

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/JrMarcco/jotice/internal/domain"
	"github.com/JrMarcco/jotice/internal/errs"
	"github.com/JrMarcco/jotice/internal/service/sendstrategy"
	"github.com/sony/sonyflake"
)

//go:generate mockgen -source=./types.go -destination=./mock/send_service.mock.go -package=notificationmock -type=SendService
type SendService interface {
	Send(ctx context.Context, n domain.Notification) (domain.SendResp, error)
	AsyncSend(ctx context.Context, n domain.Notification) (domain.SendResp, error)
	BatchSend(ctx context.Context, ns []domain.Notification) (domain.BatchSendResp, error)
	BatchAsyncSend(ctx context.Context, ns []domain.Notification) (domain.BatchAsyncSendResp, error)
}

var _ SendService = (*DefaultSendService)(nil)

type DefaultSendService struct {
	idGenerator *sonyflake.Sonyflake

	sendStrategy sendstrategy.SendStrategy
}

// Send sync send notification immediately
func (s *DefaultSendService) Send(ctx context.Context, n domain.Notification) (domain.SendResp, error) {
	resp := domain.SendResp{
		Result: domain.SendResult{
			Status: domain.SendStatusFailed,
		},
	}

	if err := n.Validate(); err != nil {
		return resp, err
	}

	id, err := s.idGenerator.NextID()
	if err != nil {
		return resp, fmt.Errorf("failed to generate notification id, cause of: %w", err)
	}

	n.Id = id

	sendResp, err := s.sendStrategy.Send(ctx, n)
	if err != nil {
		return resp, fmt.Errorf("%w, cause of: %w", errs.ErrSendNotificationFailed, err)
	}

	return sendResp, nil
}

// AsyncSend async send notification
func (s *DefaultSendService) AsyncSend(ctx context.Context, n domain.Notification) (domain.SendResp, error) {
	if err := n.Validate(); err != nil {
		return domain.SendResp{}, err
	}

	id, err := s.idGenerator.NextID()
	if err != nil {
		return domain.SendResp{}, fmt.Errorf("failed to generate notification id, cause of: %w", err)
	}

	n.Id = id

	// if immediate strategy in async send method,
	// replace strategy to deadline and set the deadline to 1 minute from now.
	n.ReplaceAsyncImmediate()
	return s.sendStrategy.Send(ctx, n)
}

// BatchSend batch send notifications
func (s *DefaultSendService) BatchSend(ctx context.Context, ns []domain.Notification) (domain.BatchSendResp, error) {
	resp := domain.BatchSendResp{}

	if len(ns) == 0 {
		return resp, fmt.Errorf("%w: no notifications to send", errs.ErrInvalidParam)
	}

	for _, n := range ns {
		if err := n.Validate(); err != nil {
			return resp, fmt.Errorf("%w: notification validation failed, cause of: %w", errs.ErrInvalidParam, err)
		}

		id, err := s.idGenerator.NextID()
		if err != nil {
			return resp, fmt.Errorf("failed to generate notification id, cause of: %w", err)
		}

		n.Id = id
	}

	results, err := s.sendStrategy.BatchSend(ctx, ns)
	resp.Results = results.Results
	if err != nil {
		return resp, fmt.Errorf("%w, cause of: %w", errs.ErrSendNotificationFailed, err)
	}

	return resp, nil
}

// BatchAsyncSend batch async send notifications
func (s *DefaultSendService) BatchAsyncSend(ctx context.Context, ns []domain.Notification) (domain.BatchAsyncSendResp, error) {
	if len(ns) == 0 {
		return domain.BatchAsyncSendResp{}, fmt.Errorf("%w: no notifications to send", errs.ErrInvalidParam)
	}

	ids := make([]uint64, 0, len(ns))
	for _, n := range ns {
		if err := n.Validate(); err != nil {
			return domain.BatchAsyncSendResp{}, fmt.Errorf("%w: notification validation failed, cause of: %w", errs.ErrInvalidParam, err)
		}

		id, err := s.idGenerator.NextID()
		if err != nil {
			return domain.BatchAsyncSendResp{}, fmt.Errorf("failed to generate notification id, cause of: %w", err)
		}

		n.Id = id
		ids = append(ids, id)
		n.ReplaceAsyncImmediate()
	}

	// group notifications by strategy
	strategyGroups := make(map[string][]domain.Notification)
	for _, n := range ns {
		strategy := string(n.StrategyConfig.Type)
		strategyGroups[strategy] = append(strategyGroups[strategy], n)
	}

	// Process each strategy group concurrently
	eg, ctx := errgroup.WithContext(ctx)
	for _, groupNs := range strategyGroups {
		notifications := groupNs
		eg.Go(func() error {
			_, err := s.sendStrategy.BatchSend(ctx, notifications)
			if err != nil {
				return fmt.Errorf("%w, cause of: %w", errs.ErrSendNotificationFailed, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return domain.BatchAsyncSendResp{}, err
	}

	return domain.BatchAsyncSendResp{
		NotificationIds: ids,
	}, nil
}

func NewDefaultSendService(idGenerator *sonyflake.Sonyflake) *DefaultSendService {
	return &DefaultSendService{
		idGenerator: idGenerator,
	}
}
