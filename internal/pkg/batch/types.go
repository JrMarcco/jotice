package batch

import (
	"context"
	"time"
)

type Adjuster interface {
	Adjust(ctx context.Context, respTime time.Duration) (int, error)
}
