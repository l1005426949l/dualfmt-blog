//go:build wireinject
// +build wireinject

package main

import (
	"my-project/metadata/biz"
	"my-project/metadata/data"
	"my-project/metadata/internal/conf"
	"my-project/metadata/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		NewGRPCServer,
		newApp,
	))
}
