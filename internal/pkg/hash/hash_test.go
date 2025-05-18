package hash

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"testing"
)

func TestHash_NoCollision(t *testing.T) {
	t.Parallel()

	testSize := 10000

	hashRes := make(map[int64]struct{}, testSize)
	inputs := make([]struct {
		bizId  int64
		bizKey string
	}, testSize)

	for i := 0; i < testSize; i++ {
		maxBig := big.NewInt(10000)
		randBig, err := rand.Int(rand.Reader, maxBig)
		if err != nil {
			t.Fatalf("failed to generate random number: %v", err)
		}
		bizId := randBig.Int64() + 1

		lenBig, err := rand.Int(rand.Reader, big.NewInt(20))
		if err != nil {
			t.Fatalf("failed to generate random number: %v", err)
		}

		keyLen := int(lenBig.Int64()) + 10
		bizKey := randomString(keyLen)

		inputs[i] = struct {
			bizId  int64
			bizKey string
		}{bizId: bizId, bizKey: bizKey}

		hashVal, err := Hash(bizId, bizKey)
		if err != nil {
			t.Fatalf("failed to hash: %v", err)
		}

		if _, exists := hashRes[hashVal]; exists {
			// collision, find inputs that have the same hash value
			for j := 0; j < i; j++ {
				v, err := Hash(inputs[j].bizId, inputs[j].bizKey)
				if err != nil {
					t.Fatalf("failed to hash: %v", err)
				}
				if v == hashVal {
					t.Fatalf("collision: %v and %v hash the same hash value: %d", inputs[j], inputs[i], hashVal)
				}
			}
		}
		hashRes[hashVal] = struct{}{}
	}

	if len(hashRes) != testSize {
		t.Fatalf("hash collision, expired number of hash: %d and actual get: %d", testSize, len(hashRes))
		return
	}

	t.Logf("success to generate %d hash values without collision", testSize)
}

//goland:noinspection SpellCheckingInspection
func randomString(l int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	res := make([]byte, l)

	// generate random bytes
	randomBytes := make([]byte, l)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("failed to generate random bytes: " + err.Error())
	}

	for i := 0; i < l; i++ {
		index := int(randomBytes[i]) % len(charset)
		res[i] = charset[index]
	}
	return string(res)
}

func TestHash_Distribution(t *testing.T) {
	t.Parallel()

	testSize := 1000000
	bucketCnt := 1000
	buckets := make([]int64, bucketCnt)

	// generate hash
	for i := 0; i < testSize; i++ {
		// rand bizId and bizKey
		maxBig := big.NewInt(10000)
		randBig, err := rand.Int(rand.Reader, maxBig)
		if err != nil {
			t.Fatalf("failed to generate random number: %v", err)
		}
		bizId := randBig.Int64() + 1
		bizKey := "test_biz_key:" + strconv.Itoa(i)

		hashVal, err := Hash(bizId, bizKey)
		if err != nil {
			t.Fatalf("failed to hash: %v", err)
		}

		// put hash value into buckets
		bucketIndex := int(hashVal % int64(bucketCnt))
		if bucketIndex < 0 {
			// deal with negative bucket index
			bucketIndex += bucketCnt
		}
		buckets[bucketIndex]++
	}

	expected := float64(testSize) / float64(bucketCnt)
	// allow 20% deviation
	maxDeviation := 0.2 * expected

	for i, cnt := range buckets {
		deviation := float64(cnt) - expected
		if deviation < 0 {
			deviation = -deviation
		}

		if deviation > maxDeviation {
			t.Logf("bucket %d hash value count: %d, deviation: %.2f, exceeding expected: %.2f", i, cnt, deviation, expected)
		}
	}

	minCnt, maxCnt, avgCnt := buckets[0], buckets[0], float64(0)
	for _, cnt := range buckets {
		if cnt < minCnt {
			minCnt = cnt
		}
		if cnt > maxCnt {
			maxCnt = cnt
		}
		avgCnt += float64(cnt)
	}
	avgCnt /= float64(bucketCnt)

	t.Logf("min hash value count: %d, max hash value count: %d, avg hash value count: %.2f", minCnt, maxCnt, avgCnt)
}
