package ring

import (
	"sync"
	"time"
)

// TimeDurationRingBuffer a fixed-size and thread-safe ring buffer of time.Duration.
type TimeDurationRingBuffer struct {
	mu sync.RWMutex

	buffer   []time.Duration
	size     int
	writePos int
	count    int
	sum      time.Duration
}

func (t *TimeDurationRingBuffer) Add(d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// buffer is full
	if t.count == t.size {
		// remove the oldest element value
		t.sum -= t.buffer[t.writePos]
	} else {
		t.count++
	}

	t.buffer[t.writePos] = d
	t.sum += d
	t.writePos = (t.writePos + 1) % t.size
}

func (t *TimeDurationRingBuffer) Avg() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.count == 0 {
		return 0
	}
	return t.sum / time.Duration(t.count)
}

func (t *TimeDurationRingBuffer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i := 0; i < t.size; i++ {
		t.buffer[i] = 0
	}
	t.count = 0
	t.sum = 0
	t.writePos = 0
}

func (t *TimeDurationRingBuffer) Size() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.size
}

func (t *TimeDurationRingBuffer) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.count
}

func NewTimeDurationRingBuffer(size int) *TimeDurationRingBuffer {
	if size <= 0 {
		panic("size must be greater than 0")
	}

	return &TimeDurationRingBuffer{
		buffer: make([]time.Duration, size),
		size:   size,
	}
}
