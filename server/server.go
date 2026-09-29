package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	meta "github.com/2manyvcos/paranal"
	apischema "github.com/2manyvcos/paranal/api"
	"github.com/2manyvcos/paranal/client"
	"github.com/2manyvcos/paranal/server/api"
	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/helper"
)

const API_PATH = "/api"

func Run() {
	log.Printf("Welcome to %s (version %s)\n", meta.Meta.Name, meta.Meta.Version)

	app, err := application.New()
	if err != nil {
		log.Fatalf("Setup failure - %s\n", err)
	}
	defer app.Close()

	apiHandler := api.New()
	http.Handle(API_PATH+"/", http.StripPrefix(API_PATH, helper.OmitTrailingSlash(helper.WithApp(app, helper.WithCORS(apiHandler)))))

	http.Handle("/", http.FileServer(helper.FileRewrite(helper.FileTemplates(client.ClientFiles, map[string]any{
		"appName": app.Config.AppName,
		"api":     app.Config.PublicAPI,
	}, "index.html", "manifest.json"), "index.html")))

	http.Handle("/api/reference/", http.FileServer(helper.FileTemplates(client.ClientFiles, map[string]any{
		"api": app.Config.PublicAPI,
	}, "api/reference/index.html")))

	apiSchemaFiles, _ := fs.Glob(apischema.SchemaFiles, "*.yaml")
	apiSchemaServer := helper.WithApp(app, helper.WithCORS(http.FileServer(helper.FileTemplates(apischema.SchemaFiles, map[string]any{
		"api": app.Config.PublicAPI,
	}, apiSchemaFiles...))))
	http.Handle("/api/schema/", http.StripPrefix("/api/schema/", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if ext := filepath.Ext(req.URL.Path); ext == ".yaml" || ext == ".yml" {
			res.Header().Set("Content-Type", "text/yaml")
		}
		apiSchemaServer.ServeHTTP(res, req)
	})))

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
