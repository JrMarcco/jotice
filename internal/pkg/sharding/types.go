package sharding

type Dst struct {
	DBPrefix    string
	TablePrefix string

	DB    string
	Table string
}

type dstContextKey struct{}

type Strategy interface {
	Shard(bizId uint64, bizKey string) Dst
	ShardWithId(id uint64) Dst
	BroadCast() []Dst
}
