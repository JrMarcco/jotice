package id

import (
	"strconv"
	"testing"

	"github.com/JrMarcco/jotice/internal/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_NextId(t *testing.T) {
	t.Parallel()

	g := NewGenerator()

	bizId := int64(2025)
	bizKey := "test_biz_key:" + strconv.FormatInt(bizId, 10)

	id, err := g.NextId(bizId, bizKey)
	require.NoError(t, err)

	hashVal := ExtractHash(id)
	if hashVal < 0 || hashVal >= 1024 {
		t.Errorf("hashVal should be in range [0, 1024), but got %d", hashVal)
	}
	wantHashVal, err := hash.Hash(bizId, bizKey)
	require.NoError(t, err)
	assert.Equal(t, hashVal, wantHashVal%1024)

	seq := ExtractSequence(id)
	if seq < 0 || seq >= (1<<12) {
		t.Errorf("seq should be in range [0, 4096), but got %d", seq)
	}
	assert.Equal(t, seq, int64(0))

}

func TestGenerator_NextId_Uniqueness(t *testing.T) {
	t.Parallel()
	g := NewGenerator()

	idCnt := 10000
	ids := make(map[int64]struct{}, idCnt)

	for i := 0; i < idCnt; i++ {
		bizId := int64(i % 100)                            // reuse bizId to test uniqueness
		bizKey := "test_biz_key:" + string(rune('A'+i%26)) // reuse bizKey to test uniqueness

		id, err := g.NextId(bizId, bizKey)
		require.NoError(t, err)

		if _, exists := ids[id]; exists {
			t.Logf("id %d already exists", id)
		}

		ids[id] = struct{}{}
	}

	conflictCnt := idCnt - len(ids)

	t.Logf("generated %d ids", idCnt)
	t.Logf("%d id conflicts and %.2f%% of the time. ",
		conflictCnt,
		float64(conflictCnt)/float64(idCnt)*float64(100),
	)
}

func TestGenerator_NextId_SeqIncr(t *testing.T) {
	t.Parallel()
	g := NewGenerator()

	biz := int64(2025)
	bizKey := "test_biz_key:" + strconv.FormatInt(biz, 10)

	cnt := 100
	ids := make([]int64, 0, cnt)
	for i := 0; i < cnt; i++ {
		id, err := g.NextId(biz, bizKey)
		require.NoError(t, err)

		seq := ExtractSequence(id)
		assert.Equal(t, seq, int64(i))

		ids = append(ids, id)
	}

	wantHashVal := ExtractHash(ids[0])
	for i := 1; i < cnt; i++ {
		hashVal := ExtractHash(ids[i])
		assert.Equal(t, hashVal, wantHashVal)
	}
}

func TestGenerator_NextId_SeqRollover(t *testing.T) {
	t.Parallel()
	g := NewGenerator()

	// change seq to max value
	g.sequence = sequenceMask

	bizId := int64(2025)
	bizKey := "test_biz_key:" + strconv.FormatInt(bizId, 10)

	id, err := g.NextId(bizId, bizKey)
	require.NoError(t, err)

	seq := ExtractSequence(id)
	assert.Equal(t, seq, int64(sequenceMask), "sequence should be max value")

	id, err = g.NextId(bizId, bizKey)
	seq = ExtractSequence(id)
	assert.Equal(t, seq, int64(0), "sequence should be rolled over to 0")
}
