package ioc

import (
	"github.com/JrMarcco/jotice/internal/pkg/snowflake"
	"go.uber.org/fx"
)

var IdGeneratorFxOpt = fx.Provide(InitIdGenerator)

func InitIdGenerator() *snowflake.Generator {
	return snowflake.NewGenerator()
}
