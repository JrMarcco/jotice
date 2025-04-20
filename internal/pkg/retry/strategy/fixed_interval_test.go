package strategy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedIntervalRetry_Next(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		strategy *FixedIntervalRetry
		wantDur  time.Duration
		wantRes  bool
	}{
		{
			name: "basic",
			strategy: &FixedIntervalRetry{
				interval:      time.Second,
				maxRetryTimes: 3,
				retriedTimes:  0,
			},
			wantDur: time.Second,
			wantRes: true,
		}, {
			name: "over max retry time",
			strategy: &FixedIntervalRetry{
				interval:      time.Second,
				maxRetryTimes: 3,
				retriedTimes:  3,
			},
			wantDur: 0,
			wantRes: false,
		}, {
			name: "max retry time is 0",
			strategy: &FixedIntervalRetry{
				interval:      time.Second,
				maxRetryTimes: 0,
				retriedTimes:  3,
			},
			wantDur: time.Second,
			wantRes: true,
		}, {
			name: "max retry time is negative",
			strategy: &FixedIntervalRetry{
				interval:      time.Second,
				maxRetryTimes: -1,
				retriedTimes:  3,
			},
			wantDur: time.Second,
			wantRes: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dur, ok := tc.strategy.Next()
			assert.Equal(t, tc.wantRes, ok)
			if ok {
				assert.Equal(t, tc.wantDur, dur)
			}
		})
	}
}

func TestFixedIntervalRetry_NextWithRetriedTimes(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		strategy     *FixedIntervalRetry
		retriedTimes int32
		wantDur      time.Duration
		wantRes      bool
	}{
		{
			name:         "basic",
			strategy:     NewFixedIntervalRetry(time.Second, 0),
			retriedTimes: 0,
			wantDur:      time.Second,
			wantRes:      true,
		}, {
			name:         "equal max retry time",
			strategy:     NewFixedIntervalRetry(time.Second, 3),
			retriedTimes: 3,
			wantDur:      time.Second,
			wantRes:      true,
		}, {
			name:         "over max retry time",
			strategy:     NewFixedIntervalRetry(time.Second, 3),
			retriedTimes: 4,
			wantDur:      0,
			wantRes:      false,
		}, {
			name:         "max retry time is 0",
			strategy:     NewFixedIntervalRetry(time.Second, 0),
			retriedTimes: 3,
			wantDur:      time.Second,
			wantRes:      true,
		}, {
			name:         "max retry time is negative",
			strategy:     NewFixedIntervalRetry(time.Second, -1),
			retriedTimes: 3,
			wantDur:      time.Second,
			wantRes:      true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dur, ok := tc.strategy.NextWithRetriedTimes(tc.retriedTimes)
			assert.Equal(t, tc.wantRes, ok)
			if ok {
				assert.Equal(t, tc.wantDur, dur)
			}
		})
	}
}
