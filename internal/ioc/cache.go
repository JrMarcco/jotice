package ioc

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var GoCacheFxOpt = fx.Provide(InitGoCache)

func InitGoCache() *cache.Cache {
	type Config struct {
		DefaultExpiration time.Duration `yaml:"defaultExpiration"`
		CleanupInterval   time.Duration `yaml:"cleanupInterval"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("cache", &cfg); err != nil {
		panic(err)
	}

	return cache.New(cfg.DefaultExpiration, cfg.CleanupInterval)
}
