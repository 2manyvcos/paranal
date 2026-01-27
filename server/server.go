package server

import (
	"log"

	meta "github.com/2manyvcos/paranal"
	"github.com/2manyvcos/paranal/server/api"
	"github.com/2manyvcos/paranal/server/application"
)

func Run() {
	log.Printf("Welcome to %s v%s\n", meta.Meta.Name, meta.Meta.Version)

	app, err := application.Setup()
	if err != nil {
		log.Fatalf("Setup failure - %s\n", err)
	}
	defer app.Close()

	err = api.Serve(app)
	if err != nil {
		log.Printf("HTTP server exited - %s\n", err)
	}
}
