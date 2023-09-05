//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"synchydra/internal/pkg/middleware"
	"synchydra/internal/repository"
	"synchydra/pkg/log"
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	middleware.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
)

func newApp(*viper.Viper, *log.Logger) (*Migrate, func(), error) {
	panic(wire.Build(
		RepositorySet,
		NewMigrate,
	))
}
