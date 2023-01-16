// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	"github.com/frchandra/gmcgo/app"
	"github.com/frchandra/gmcgo/config"
	"github.com/frchandra/gmcgo/database"
)

// Injectors from injector.go:

func InitializeServer() *app.Server {
	appConfig := config.NewAppConfig()
	server := app.NewServer(appConfig)
	return server
}

func InitializeMigrator() *database.Migrator {
	appConfig := config.NewAppConfig()
	migration := database.NewMigration()
	migrator := database.NewMigrator(appConfig, migration)
	return migrator
}
