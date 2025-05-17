package sharding

var _ Strategy = (*HashStrategy)(nil)

type HashStrategy struct {
	dbPrefix    string
	tablePrefix string

	dbSharding    int32
	tableSharding int32
}

func (h HashStrategy) Shard(bizId uint64, bizKey string) Dst {
	//TODO implement me
	panic("implement me")
}

func (h HashStrategy) ShardWithId(id uint64) Dst {
	//TODO implement me
	panic("implement me")
}

func (h HashStrategy) BroadCast() []Dst {
	//TODO implement me
	panic("implement me")
}

func NewHashStrategy(dbPrefix, tablePrefix string, dbSharding, tableSharding int32) HashStrategy {
	return HashStrategy{
		dbPrefix:      dbPrefix,
		tablePrefix:   tablePrefix,
		dbSharding:    dbSharding,
		tableSharding: tableSharding,
	}
}
