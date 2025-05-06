package domain

import (
	"fmt"
	"time"

	"github.com/JrMarcco/jotice/internal/errs"
)

type SendStrategy string

const (
	SendStrategyImmediate  SendStrategy = "immediate"
	SendStrategyDelayed    SendStrategy = "delayed"
	SendStrategyScheduled  SendStrategy = "scheduled"
	SendStrategyTimeWindow SendStrategy = "time_window"
	SendStrategyDeadline   SendStrategy = "deadline"
)

type SendStrategyConfig struct {
	Strategy   SendStrategy
	Delay      time.Duration
	ScheduleAt time.Time
	Start      time.Time
	End        time.Time
	Deadline   time.Time
}

// Validate validate the strategy config
func (c SendStrategyConfig) Validate() error {
	switch c.Strategy {
	case SendStrategyImmediate:
		return nil
	case SendStrategyDelayed:
		if c.Delay <= 0 {
			return fmt.Errorf("%w: delay time should be greater than 0", errs.ErrInvalidParam)
		}
	case SendStrategyScheduled:
		if c.ScheduleAt.IsZero() || c.ScheduleAt.Before(time.Now()) {
			return fmt.Errorf("%w: schedule_at should not be zero or before now", errs.ErrInvalidParam)
		}
	case SendStrategyTimeWindow:
		if c.Start.IsZero() || c.Start.After(c.End) {
			return fmt.Errorf("%w: start and end time should not be zero and start should be before end", errs.ErrInvalidParam)
		}
	case SendStrategyDeadline:
		if c.Deadline.IsZero() || c.Deadline.Before(time.Now()) {
			return fmt.Errorf("%w: deadline should not be zero or before now", errs.ErrInvalidParam)
		}
	default:
		return fmt.Errorf("%w: unknown strategy", errs.ErrInvalidParam)
	}

	return nil
}

// CalcTimeWindow calculate the start and end time based on the strategy to send the notification
func (c SendStrategyConfig) CalcTimeWindow() (start, end time.Time) {
	switch c.Strategy {
	case SendStrategyImmediate:
		//
		now := time.Now()
		return now, now
	case SendStrategyDelayed:
		return time.Now(), time.Now().Add(c.Delay)
	case SendStrategyScheduled:
		return c.ScheduleAt, c.ScheduleAt
	case SendStrategyTimeWindow:
		return c.Start, c.End
	case SendStrategyDeadline:
		return c.Start.Add(-3 * time.Second), c.Deadline
	default:
		now := time.Now()
		return now, now
	}
}

type SendResult struct {
	NotificationId uint64
	Status         SendStatus
}

type SendResp struct {
	Result SendResult
}

type BatchSendResp struct {
	Results []SendResult
}

type BatchAsyncSendResp struct {
	NotificationIds []uint64
}
