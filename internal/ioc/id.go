package ioc

import (
	"github.com/JrMarcco/jotice/internal/pkg/id"
	"go.uber.org/fx"
)

var IdGeneratorFxOpt = fx.Provide(InitIdGenerator)

func InitIdGenerator() *id.Generator {
	return id.NewGenerator()
}
