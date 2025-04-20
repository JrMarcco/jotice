package strategy

import (
	"sync/atomic"
	"time"
)

var _ Strategy = (*FixedIntervalRetry)(nil)

type FixedIntervalRetry struct {
	interval      time.Duration
	maxRetryTimes int32
	retriedTimes  int32
}

func (f *FixedIntervalRetry) Next() (time.Duration, bool) {
	retriedTimes := atomic.AddInt32(&f.retriedTimes, 1)
	return f.nextRetry(retriedTimes)
}

func (f *FixedIntervalRetry) NextWithRetriedTimes(retriedTimes int32) (time.Duration, bool) {
	return f.nextRetry(retriedTimes)
}

func (f *FixedIntervalRetry) nextRetry(retriedTimes int32) (time.Duration, bool) {
	if f.maxRetryTimes <= 0 || retriedTimes <= f.maxRetryTimes {
		return f.interval, true
	}
	return 0, false
}

func NewFixedIntervalRetry(interval time.Duration, maxRetryTime int32) *FixedIntervalRetry {
	return &FixedIntervalRetry{
		interval:      interval,
		maxRetryTimes: maxRetryTime,
		retriedTimes:  0,
	}
}
