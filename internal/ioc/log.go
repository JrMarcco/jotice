package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var LoggerFxOpt = fx.Provide(
	fx.Annotate(
		InitLogger,
		fx.ResultTags(`name:"zapLogger"`),
	),
)

func InitLogger() *zap.Logger {
	type config struct {
		Env string `yaml:"env"`
	}

	cfg := config{
		Env: "dev",
	}

	err := viper.UnmarshalKey("profile", &cfg)
	if err != nil {
		panic(err)
	}

	var zl *zap.Logger
	switch cfg.Env {
	case "prod":
		zl, err = zap.NewProduction()
	default:
		zl, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	return zl
}
