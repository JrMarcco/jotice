package id

import (
	"sync/atomic"
	"time"

	"github.com/JrMarcco/jotice/internal/pkg/hash"
)

const (
	timestampBits = 41
	hashBits      = 10
	sequenceBits  = 12

	hashShift      = sequenceBits
	timestampShift = hashShift + hashBits

	sequenceMask  = (1 << sequenceBits) - 1
	hashMask      = (1 << hashBits) - 1
	timestampMask = (1 << timestampBits) - 1

	epochMillis   = int64(1735689600000) // milliseconds of 2025-01-01 00:00:00
	number1000    = int64(1000)
	number1000000 = int64(1000000)
)

type Generator struct {
	sequence int64
	lastTime int64     // the time of last id generated
	epoch    time.Time // the epoch time
}

func NewGenerator() *Generator {
	return &Generator{
		sequence: 0,
		lastTime: 0,
		epoch:    time.Unix(epochMillis/number1000, (epochMillis%number1000)*number1000000),
	}
}

// NextId returns the next id.
//
// The id is composed of:
//   - 41 bits for timestamp, the timestamp is the milliseconds of the current time minus the epoch time.
//   - 10 bits for hash, the hash is the hash value of the bizId and bizKey.
//     bizId and bizKey decide the database sharding.
//   - 12 bits for a sequence, the sequence is the auto incr sequence number of the id.
func (g *Generator) NextId(bizId int64, bizKey string) (int64, error) {
	timestamp := time.Now().UnixMilli() - epochMillis
	hashVal, err := hash.Hash(bizId, bizKey)
	if err != nil {
		return 0, err
	}
	seq := atomic.AddInt64(&g.sequence, 1) - 1

	nextId := (timestamp&timestampMask)<<timestampShift | (hashVal&hashMask)<<hashShift | (seq & sequenceMask)
	return nextId, nil
}

func ExtractTimestamp(id int64) time.Time {
	timestamp := (id >> timestampShift) % timestampMask
	return time.Unix(0, (timestamp+epochMillis)*int64(time.Millisecond))
}

func ExtractHash(id int64) int64 {
	return (id >> hashShift) & hashMask
}

func ExtractSequence(id int64) int64 {
	return id & sequenceMask
}
