package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	meta "github.com/2manyvcos/paranal"
	"github.com/2manyvcos/paranal/client"
	"github.com/2manyvcos/paranal/server/api"
	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/state"
)

const API_PATH = "/api"

func Run() {
	log.Printf("Welcome to %s v%s\n", meta.Meta.Name, meta.Meta.Version)

	app, err := application.Setup()
	if err != nil {
		log.Fatalf("Setup failure - %s\n", err)
	}
	defer app.Close()

	app.State, err = state.Setup(app)
	if err != nil {
		app.Close()
		log.Fatalf("Setup failure - %s\n", err)
	}

	app.Scheduler.Start()
	app.RunServiceScripts()

	apiHandler := api.New()
	http.Handle(API_PATH+"/", http.StripPrefix(API_PATH, helper.OmitTrailingSlash(helper.WithApp(app, helper.WithCORS(apiHandler)))))

	resolvedClientFiles := helper.FileTemplates(client.ClientFiles, map[string]any{
		"appName":           app.Config.AppName,
		"tagline":           app.Config.Tagline,
		"logo":              app.Config.Brand.Logo,
		"favicon":           app.Config.Brand.Favicon,
		"appleTouchIcon":    app.Config.Brand.AppleTouchIcon,
		"theme":             app.Config.Brand.Theme,
		"logoutRedirectURL": app.Config.Auth.LogoutRedirectURL,
	}, "index.html", "manifest.json")
	http.Handle("/", http.FileServer(helper.FileRewrite(resolvedClientFiles, "index.html")))

	hostname := fmt.Sprintf("%s:%s", app.Config.Server.Address, app.Config.Server.Port)
	switch strings.ToLower(app.Config.Server.Protocol) {
	case "http":
		log.Printf("Listening at http://%s\n", hostname)
		err = http.ListenAndServe(hostname, nil)
		if err != nil {
			log.Printf("HTTP server exited - %s\n", err)
		}
	case "https":
		log.Printf("Listening at https://%s\n", hostname)
		err = http.ListenAndServeTLS(hostname, app.Config.Server.CertFile, app.Config.Server.KeyFile, nil)
		if err != nil {
			log.Printf("HTTP server exited - %s\n", err)
		}
	default:
		log.Printf("Error starting HTTP server - unsupported protocol \"%s\"\n", app.Config.Server.Protocol)
	}
}
