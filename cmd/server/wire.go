//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"synchydra/internal/handler"
	"synchydra/internal/repository"
	"synchydra/internal/server"
	"synchydra/internal/service"
	"synchydra/pkg/helper/sid"
	"synchydra/pkg/log"
)

var middlewareSet = wire.NewSet(
	middleware.NewRedis,
	middleware.NewCanalClient,
	middleware.NewRocketmqProvider,
	middleware.NewRocketmqConsumer,
)epository.NewUserRepository,
)

func newApp(*viper.Viper, *log.Logger) (*server.Server, func(), error) {
	panic(wire.Build(
		middlewareSet,
		server.NewServer,
		server.NewServerHTTP,
		sid.NewSid,
	))
}
