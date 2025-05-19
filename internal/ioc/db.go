package ioc

import (
	"context"
	"database/sql"
	"time"

	"github.com/JrMarcco/easy-kit/retry"
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
		DSN: "root:root@tcp(localhost:3306)/jotice?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=1s&readTimeout=3s&writeTimeout=3s&multiStatements=true&interpolateParams=true",
	}

	if err := viper.UnmarshalKey("db.postgres", &cfg); err != nil {
		panic(err)
	}

	waitForDBSetup(cfg.DSN)

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func waitForDBSetup(dsn string) {
	sqlDB, err := sql.Open("postgres", dsn)
	defer func(sqlDB *sql.DB) { _ = sqlDB.Close() }(sqlDB)
	if err != nil {
		panic(err)
	}

	const maxInterval = 10 * time.Second
	const maxRetries = 10

	strategy, err := retry.NewExponentialBackoffStrategy(time.Second, maxInterval, maxRetries)
	if err != nil {
		panic(err)
	}

	const timeout = time.Second
	for {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		err = sqlDB.PingContext(ctx)
		cancel()

		if err == nil {
			break
		}

		next, ok := strategy.Next()
		if !ok {
			panic("failed to connect to database after max retries: " + err.Error())
		}
		time.Sleep(next)
	}
}
