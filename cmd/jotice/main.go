package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func main() {
	// init config
	initViper()

	fx.New().Run()
}

func initViper() {
	configFile := pflag.String("config", "config/config.yaml", "Specify config file path")
	pflag.Parse()

	viper.SetConfigFile(*configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
