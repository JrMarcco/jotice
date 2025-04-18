package domain

import (
	"fmt"
	"time"

	"github.com/JrMarcco/jotice/internal/errs"
)

type SendStatus string

const (
	SendStatusPrepare  SendStatus = "prepare"
	SendStatusCanceled SendStatus = "canceled"
	SendStatusPending  SendStatus = "pending"
	SendStatusSending  SendStatus = "sending"
	SendStatusSuccess  SendStatus = "success"
	SendStatusFailed   SendStatus = "failed"
)

func (s SendStatus) String() string {
	return string(s)
}

type Notification struct {
	Id             uint64
	BizId          int64
	Key            string
	Receivers      []string
	Channel        Channel
	Template       Template
	Status         SendStatus
	ScheduledStart time.Time
	ScheduledEnd   time.Time
	Version        int32
	StrategyConfig SendStrategyConfig
}

type Template struct {
	Id        uint64
	VersionId uint64
	Params    map[string]string
}

func (n *Notification) Validate() error {
	if n.BizId <= 0 {
		return fmt.Errorf("%w: biz_id = %q", errs.ErrInvalidParam, n.BizId)
	}

	if n.Key == "" {
		return fmt.Errorf("%w: key = %q", errs.ErrInvalidParam, n.Key)
	}

	if len(n.Receivers) == 0 {
		return fmt.Errorf("%w: receivers = %q", errs.ErrInvalidParam, n.Receivers)
	}

	if err := n.StrategyConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func (n *Notification) IsImmediate() bool {
	return n.StrategyConfig.Strategy == SendStrategyImmediate
}

// ReplaceAsyncImmediate replace the notification to be sent asynchronously
func (n *Notification) ReplaceAsyncImmediate() {
	if n.IsImmediate() {
		n.StrategyConfig.Deadline = time.Now().Add(time.Minute)
		n.StrategyConfig.Strategy = SendStrategyDeadline
	}
}
