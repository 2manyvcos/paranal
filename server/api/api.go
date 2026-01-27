package api

import (
	"log"

	"github.com/2manyvcos/paranal/server/application"
)

func Serve(app *application.App) error {
	log.Printf("HTTP server is listening at %s://%s:%s\n", app.Config.Server.Protocol, app.Config.Server.Address, app.Config.Server.Port)

	return nil
}
