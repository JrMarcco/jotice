package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBFxOpt = fx.Provide(initDB)

func initDB() *gorm.DB {
	type config struct {
		DSN string `yaml:"dsn"`
	}

	cfg := config{
		DSN: "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
	}

	err := viper.UnmarshalKey("db.postgres", &cfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
