package batch

import (
	"context"
	"time"
)

var _ Adjuster = (*FixedStepAdjuster)(nil)

type FixedStepAdjuster struct {
	batchSize    int
	minBatchSize int
	maxBatchSize int
	adjustStep   int

	lastAdjustAt   time.Time
	adjustInterval time.Duration

	fastRespThreshold time.Duration
	slowRespThreshold time.Duration
}

func (f *FixedStepAdjuster) Adjust(_ context.Context, respTime time.Duration) (int, error) {
	if !f.lastAdjustAt.IsZero() && time.Since(f.lastAdjustAt) < f.adjustInterval {
		return f.batchSize, nil
	}

	if respTime < f.fastRespThreshold {
		// response time less than the fast response time threshold, adjust the batch size
		if f.batchSize < f.maxBatchSize {
			f.batchSize = min(f.batchSize+f.adjustStep, f.maxBatchSize)
			f.lastAdjustAt = time.Now()
		}
		return f.batchSize, nil
	}

	if respTime > f.slowRespThreshold {
		// response time more than the slow response time threshold, adjust the batch size
		if f.batchSize > f.minBatchSize {
			f.batchSize = max(f.batchSize-f.adjustStep, f.minBatchSize)
			f.lastAdjustAt = time.Now()
		}
	}
	return f.batchSize, nil
}

func NewFixedStepAdjuster(
	batchSize, minBatchSize, maxBatchSize, adjustStep int,
	adjustInterval, fastRespThreshold, slowRespThreshold time.Duration,
) *FixedStepAdjuster {
	if batchSize < minBatchSize {
		batchSize = minBatchSize
	}

	if batchSize > maxBatchSize {
		batchSize = maxBatchSize
	}

	return &FixedStepAdjuster{
		batchSize:         batchSize,
		minBatchSize:      minBatchSize,
		maxBatchSize:      maxBatchSize,
		adjustStep:        adjustStep,
		adjustInterval:    adjustInterval,
		fastRespThreshold: fastRespThreshold,
		slowRespThreshold: slowRespThreshold,
		lastAdjustAt:      time.Time{},
	}
}
