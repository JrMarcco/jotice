package retry

import (
	"fmt"
	"github.com/JrMarcco/jotice/internal/pkg/retry/strategy"
	"time"
)

type Config struct {
	Type               string                    `json:"type"`
	FixedInterval      *FixedIntervalConfig      `json:"fixed_interval"`
	ExponentialBackoff *ExponentialBackoffConfig `json:"exponential_backoff"`
}

type ExponentialBackoffConfig struct {
	InitialInterval time.Duration `json:"initial_interval"`
	MaxInterval     time.Duration `json:"max_interval"`
	MaxRetryTime    int32         `json:"max_retry_time"`
}

type FixedIntervalConfig struct {
	Interval     time.Duration `json:"interval"`
	MaxRetryTime int32         `json:"max_retry_time"`
}

func NewRetryStrategy(cfg Config) (strategy.Strategy, error) {
	switch cfg.Type {
	case "fixed_interval":
		return strategy.NewFixedIntervalRetry(cfg.FixedInterval.Interval, cfg.FixedInterval.MaxRetryTime), nil
	case "exponential_backoff":
		return strategy.NewExponentialBackoffRetry(
			cfg.ExponentialBackoff.InitialInterval,
			cfg.ExponentialBackoff.MaxInterval,
			cfg.ExponentialBackoff.MaxRetryTime,
		), nil
	default:
		return nil, fmt.Errorf("unknown retry strategy type: %s", cfg.Type)
	}
}
