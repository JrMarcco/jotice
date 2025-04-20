package strategy

import (
	"math"
	"sync/atomic"
	"time"
)

var _ Strategy = (*ExponentialBackoffRetry)(nil)

type ExponentialBackoffRetry struct {
	initialInterval    time.Duration
	maxInterval        time.Duration
	maxRetryTimes      int32
	retriedTimes       int32
	reachedMaxInterval atomic.Bool
}

func (e *ExponentialBackoffRetry) Next() (time.Duration, bool) {
	retriedTimes := atomic.AddInt32(&e.retriedTimes, 1)
	return e.nextRetry(retriedTimes)
}

func (e *ExponentialBackoffRetry) NextWithRetriedTimes(retriedTimes int32) (time.Duration, bool) {
	return e.nextRetry(retriedTimes)
}

func (e *ExponentialBackoffRetry) nextRetry(retriedTimes int32) (time.Duration, bool) {
	if e.maxRetryTimes <= 0 || retriedTimes <= e.maxRetryTimes {
		if e.reachedMaxInterval.Load() {
			return e.maxInterval, true
		}

		const two = 2
		interval := e.initialInterval * time.Duration(math.Pow(two, float64(retriedTimes-1)))

		// interval = 0 prevents an input interval = 0 when create strategy.
		// interval < 0 means the interval is over max int32 value after math.Pow.
		if interval <= 0 || interval > e.maxInterval {
			e.reachedMaxInterval.Store(true)
			return e.maxInterval, true
		}

		return interval, true
	}

	return 0, false
}

func NewExponentialBackoffRetry(initialInterval, maxInterval time.Duration, maxRetryTime int32) *ExponentialBackoffRetry {
	return &ExponentialBackoffRetry{
		initialInterval: initialInterval,
		maxInterval:     maxInterval,
		maxRetryTimes:   maxRetryTime,
	}
}
