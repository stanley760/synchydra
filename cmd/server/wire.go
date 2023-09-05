//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"synchydra/internal/handler"
	"synchydra/internal/pkg/middleware"
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
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
)

func newApp(*viper.Viper, *log.Logger) (*server.Server, func(), error) {
	panic(wire.Build(
		middlewareSet,
		repositorySet,
		serviceSet,
		handlerSet,
		server.NewServer,
		server.NewServerHTTP,
		sid.NewSid,
	))
}
