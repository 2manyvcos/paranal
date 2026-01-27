package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/2manyvcos/paranal/client"
	"github.com/2manyvcos/paranal/server/application"
)

const API_PREFIX = "/api"
const API_VERSION = "v1"
const API_PATH = API_PREFIX + "/" + API_VERSION

func Serve(app *application.App) error {
	api := http.NewServeMux()
	// api.HandleFunc("POST /auth", HandleAuth(p))
	// api.HandleFunc("GET /demo-endpoint", AuthValidator(p, HandleDemoEndpoint(p)))
	api.HandleFunc("GET /echo/{text}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.PathValue("text")))
	})

	http.Handle(API_PATH+"/", http.StripPrefix(API_PATH, OmitTrailingSlash(api)))
	http.HandleFunc(API_PREFIX+"/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	http.Handle("/", http.FileServer(FileRewrite(http.FS(client.ClientFiles), "index.html")))

	hostname := fmt.Sprintf("%s:%s", app.Config.Server.Address, app.Config.Server.Port)

	switch strings.ToLower(app.Config.Server.Protocol) {
	case "http":
		log.Printf("Listening at http://%s\n", hostname)
		return http.ListenAndServe(hostname, nil)

	case "https":
		log.Printf("Listening at https://%s\n", hostname)
		return http.ListenAndServeTLS(hostname, app.Config.Server.CertFile, app.Config.Server.KeyFile, nil)

	default:
		return fmt.Errorf(`unsupported protocol "%s"`, app.Config.Server.Protocol)
	}
}
