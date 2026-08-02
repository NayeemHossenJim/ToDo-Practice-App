package main

import (
	"go.uber.org/fx"

	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/config"
	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/database"
	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/server"
	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/todo"
)

func applicationOptions() []fx.Option {
	return []fx.Option{
		fx.Provide(
			config.Load,
			database.NewPostgresPool,
			database.NewQuerier,
			todo.NewHandler,
			server.NewRouter,
			server.NewHTTPServer,
		),
		fx.Invoke(
			server.RegisterLifecycle,
		),
	}
}
