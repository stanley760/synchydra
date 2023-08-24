//go:build wireinject
// +build wireinject

package main

import (
	"synchydra/internal/handler"
	"synchydra/internal/repository"
	"synchydra/internal/server"
	"synchydra/internal/service"
	"synchydra/pkg/helper/sid"
	"synchydra/pkg/jwt"
	"synchydra/pkg/log"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
)

func newApp(*viper.Viper, *log.Logger) (*server.Server, func(), error) {
	panic(wire.Build(
		RepositorySet,
		ServiceSet,
		HandlerSet,
		server.NewServer,
		server.NewServerHTTP,
		sid.NewSid,
		jwt.NewJwt,
	))
}
