package ring

import "sync"

const (
	bitsPerWord = 64              // number of bits in uint64
	bitsMask    = bitsPerWord - 1 // mask for the bit position within uint64
	bitsShift   = 6               // equivalent to division by 64 (2^6)

	defaultWindowSize     = 128 // default window size
	defaultMinConsecutive = 3   // default minimum consecutive events
)

// BitRing is a sliding window ring buffer using bits to efficiently record event occurrences.
type BitRing struct {
	mu sync.RWMutex // read-write lock for concurrency safety

	bits []uint64 // bit array to store event states

	windowSize int // length of the window (ring buffer size)
	writePos   int // next position to write an event

	isFull           bool    // whether the ring buffer is full
	eventCount       int     // number of events occurred in the current window
	minConsecutive   int     // minimum consecutive events to trigger
	triggerThreshold float64 // threshold for event occurrence rate (0~1)
}

// NewBitRing creates a new BitRing instance.
// windowSize: size of the window; minConsecutive: minimum consecutive events to trigger; threshold: event occurrence rate threshold.
func NewBitRing(windowSize int, minConsecutive int, threshold float64) *BitRing {
	if windowSize <= 0 {
		windowSize = defaultWindowSize
	}
	if minConsecutive <= 0 {
		minConsecutive = defaultMinConsecutive
	}

	if minConsecutive > windowSize {
		minConsecutive = windowSize
	}

	if threshold < 0 {
		threshold = 0
	}
	if threshold > 1 {
		threshold = 1
	}

	return &BitRing{
		bits:             make([]uint64, (windowSize+bitsMask)/bitsPerWord),
		windowSize:       windowSize,
		minConsecutive:   minConsecutive,
		triggerThreshold: threshold,
	}
}

// Add inserts an event into the sliding window.
// eventHappened: whether the event occurred (true/false).
func (b *BitRing) Add(eventHappened bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	oldBit := b.bitAt(b.writePos) // record the original event state at the current position

	// If the ring is full and the old position was 1, decrease the event count
	if b.isFull && oldBit {
		b.eventCount--
	}
	b.setBit(b.writePos, eventHappened) // write the new event
	if eventHappened {
		b.eventCount++
	}

	b.writePos++
	if b.writePos >= b.windowSize {
		b.writePos = 0
		b.isFull = true
	}
}

// bitAt gets the bit value at the specified index.
// Returns true if the bit is 1 (event occurred).
func (b *BitRing) bitAt(index int) bool {
	pos := index >> bitsShift        // index in the uint64 array
	offset := uint(index & bitsMask) // the bit offset within uint64
	return (b.bits[pos]>>offset)&1 == 1
}

// setBit sets the bit value at the specified index.
// val=true sets to 1, val=false sets to 0.
func (b *BitRing) setBit(index int, val bool) {
	pos := index >> bitsShift
	offset := uint(index & bitsMask)

	if val {
		// set to 1
		b.bits[pos] |= 1 << offset
		return
	}
	// set to 0 (using &^ for the bit clear)
	// &^ is a special bit operation in Go, using to zeroing the bit of the left operand where the right operand is 1.
	b.bits[pos] &^= 1 << offset
}

// ShouldTrigger determines whether the current window meets the trigger condition.
// Trigger conditions:
// 1. The last minConsecutive events all occurred;
// 2. Or the event occurrence rate exceeds the triggerThreshold.
func (b *BitRing) ShouldTrigger() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	currSize := b.currWindowSize()
	if currSize == 0 {
		return false
	}

	// check if the last minConsecutive events all occurred
	if currSize >= b.minConsecutive {
		all := true
		for i := 1; i <= b.minConsecutive; i++ {
			pos := (b.writePos - i + b.windowSize) % b.windowSize
			if b.bitAt(pos) {
				continue
			}
			all = false
			break
		}
		if all {
			return true
		}
	}

	// check if the event occurrence rate exceeds the threshold
	if float64(b.eventCount)/float64(currSize) > b.triggerThreshold {
		return true
	}
	return false
}

// currWindowSize returns the actual size of the current window.
// If the ring is full, returns windowSize; otherwise, returns the number of written events.
func (b *BitRing) currWindowSize() int {
	if b.isFull {
		return b.windowSize
	}
	return b.writePos
}
