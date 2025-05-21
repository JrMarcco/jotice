package idempotent

import "context"

// Strategy idempotent strategy
type Strategy interface {
	Exists(ctx context.Context, key string) (bool, error)
	MultiExists(ctx context.Context, keys []string) (map[string]bool, error)
}
