package strategy

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExponentialBackoffRetry_Next(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		strategy *ExponentialBackoffRetry
		wantDur  time.Duration
		wantRes  bool
	}{
		{
			name: "basic",
			strategy: &ExponentialBackoffRetry{
				initialInterval: time.Second,
				maxInterval:     time.Minute,
				maxRetryTimes:   3,
				retriedTimes:    0,
			},
			wantDur: time.Second,
			wantRes: true,
		}, {
			name: "over max retry time",
			strategy: &ExponentialBackoffRetry{
				initialInterval: time.Second,
				maxInterval:     time.Minute,
				maxRetryTimes:   3,
				retriedTimes:    3,
			},
			wantDur: 0,
			wantRes: false,
		}, {
			name: "initial interval is 0",
			strategy: &ExponentialBackoffRetry{
				initialInterval: 0,
				maxInterval:     time.Minute,
				maxRetryTimes:   3,
				retriedTimes:    0,
			},
			wantDur: time.Minute,
			wantRes: true,
		}, {
			name: "over max interval",
			strategy: &ExponentialBackoffRetry{
				initialInterval: time.Second,
				maxInterval:     time.Minute,
				maxRetryTimes:   10,
				retriedTimes:    8,
			},
			wantDur: time.Minute,
			wantRes: true,
		}, {
			name: "in max interval",
			strategy: &ExponentialBackoffRetry{
				initialInterval: time.Second,
				maxInterval:     time.Minute,
				maxRetryTimes:   10,
				retriedTimes:    2,
			},
			wantDur: 4 * time.Second,
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

func TestExponentialBackoffRetry_NextWithRetriedTimes(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		strategy     *ExponentialBackoffRetry
		retriedTimes int32
		wantDur      time.Duration
		wantRes      bool
	}{
		{
			name:         "basic",
			strategy:     NewExponentialBackoffRetry(time.Second, time.Minute, 3),
			retriedTimes: 1,
			wantDur:      time.Second,
			wantRes:      true,
		}, {
			name:         "equal max retry time",
			strategy:     NewExponentialBackoffRetry(time.Second, time.Minute, 3),
			retriedTimes: 3,
			wantDur:      4 * time.Second,
			wantRes:      true,
		}, {
			name:         "over max retry time",
			strategy:     NewExponentialBackoffRetry(time.Second, time.Minute, 3),
			retriedTimes: 4,
			wantDur:      0,
			wantRes:      false,
		}, {
			name:         "initial interval is 0",
			strategy:     NewExponentialBackoffRetry(0, time.Minute, 3),
			retriedTimes: 1,
			wantDur:      time.Minute,
			wantRes:      true,
		}, {
			name:         "over max interval",
			strategy:     NewExponentialBackoffRetry(time.Second, time.Minute, 10),
			retriedTimes: 8,
			wantDur:      time.Minute,
			wantRes:      true,
		}, {
			name:         "in max interval",
			strategy:     NewExponentialBackoffRetry(time.Second, time.Minute, 10),
			retriedTimes: 3,
			wantDur:      4 * time.Second,
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
