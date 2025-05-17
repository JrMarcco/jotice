package client

const (
	attrNameWeight      = "attr_weight"
	attrNameReadWeight  = "attr_read_weight"
	attrNameWriteWeight = "attr_write_weight"
	attrNameGroup       = "attr_group"
	attrNameNode        = "attr_node"
)

type groupKey struct{}

type ContextKey string

const KeyRequestType ContextKey = "request_type"
