package api

import "github.com/2manyvcos/paranal/server/config"

type App struct {
	Config *config.Config
}

func NewApp(cfg *config.Config) *App {
	var app App
	app.Config = cfg

	return &app
}
