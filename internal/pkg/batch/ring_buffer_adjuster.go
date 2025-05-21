package batch

import (
	"context"
	"sync"
	"time"

	"github.com/JrMarcco/jotice/internal/pkg/ring"
)

const (
	defaultBufferSize = 128
)

type RingBufferAdjuster struct {
	mu sync.RWMutex

	buffer *ring.TimeDurationRingBuffer

	batchSize    int
	minBatchSize int
	maxBatchSize int
	adjustStep   int

	lastAdjustAt   time.Time
	adjustInterval time.Duration
}

func (r *RingBufferAdjuster) Adjust(_ context.Context, respTime time.Duration) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.buffer.Add(respTime)

	if !r.lastAdjustAt.IsZero() && time.Since(r.lastAdjustAt) < r.adjustInterval {
		return r.batchSize, nil
	}

	// buffer is not full
	if r.buffer.Count() < r.buffer.Size() {
		return r.batchSize, nil
	}

	avgTime := r.buffer.Avg()
	if respTime < avgTime {
		// response time less than the average time, adjust the batch size
		if r.batchSize < r.maxBatchSize {
			r.batchSize = min(r.batchSize+r.adjustStep, r.maxBatchSize)
			r.lastAdjustAt = time.Now()
		}
		return r.batchSize, nil
	}

	if respTime > avgTime {
		// response time more than the average time, adjust the batch size
		if r.batchSize > r.minBatchSize {
			r.batchSize = max(r.batchSize-r.adjustStep, r.minBatchSize)
			r.lastAdjustAt = time.Now()
		}
	}

	return r.batchSize, nil
}

func NewRingBufferAdjuster(
	bufferSize, batchSize, minBatchSize, maxBatchSize, adjustStep int, adjustInterval time.Duration,
) *RingBufferAdjuster {
	if bufferSize <= 0 {
		bufferSize = defaultBufferSize
	}
	if batchSize < minBatchSize {
		batchSize = minBatchSize
	}

	if batchSize > maxBatchSize {
		batchSize = maxBatchSize
	}

	buffer := ring.NewTimeDurationRingBuffer(bufferSize)
	return &RingBufferAdjuster{
		buffer:         buffer,
		batchSize:      batchSize,
		minBatchSize:   minBatchSize,
		maxBatchSize:   maxBatchSize,
		adjustStep:     adjustStep,
		adjustInterval: adjustInterval,
		lastAdjustAt:   time.Time{},
	}
}
