//go:build wireinject
// +build wireinject

package main

import (
	"file-management-service/internal/biz"
	"file-management-service/internal/conf"
	"file-management-service/internal/data"
	"file-management-service/internal/server"
	"file-management-service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
