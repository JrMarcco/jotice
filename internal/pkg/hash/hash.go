package hash

import (
	"hash/fnv"
	"math/bits"
	"strconv"
)

const (
	hashMask = 0x7FFFFFFFFFFFFFFF

	number13 = 13
	number29 = 29
	number31 = 31

	prime1 = 11400714819323198485
	prime2 = 14029467366897019727
	prime3 = 1609587929392839161
)

func Hash(bizId int64, bizKey string) (int64, error) {
	key := strconv.FormatInt(bizId, 10) + ":" + bizKey

	hash64a := fnv.New64a()
	_, err := hash64a.Write([]byte(key))
	if err != nil {
		return 0, err
	}
	hashVal := hash64a.Sum64()
	hashVal = mixHash(hashVal, uint64(bizId))
	return int64(hashVal) & hashMask, nil
}

func mixHash(hashVal uint64, salt uint64) uint64 {
	hashVal ^= salt + prime1
	hashVal = bits.RotateLeft64(hashVal, number13)
	hashVal *= prime2
	hashVal = bits.RotateLeft64(hashVal, number29)
	hashVal *= prime3
	hashVal = bits.RotateLeft64(hashVal, number31)

	return hashVal
}
