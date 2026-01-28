package api

import (
	"net/http"

	"github.com/2manyvcos/paranal/server/helper"
)

const API_VERSION = "v1"

func New() http.Handler {
	api := http.NewServeMux()
	// api.HandleFunc("POST /auth", HandleAuth(p))
	// api.HandleFunc("GET /demo-endpoint", AuthValidator(p, HandleDemoEndpoint(p)))
	// api.HandleFunc("GET /echo/{text}", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte(r.PathValue("text")))
	// })
	api.HandleFunc("POST /auth", PostAuth)
	api.HandleFunc("GET /user", helper.WithAuth(http.HandlerFunc(GetUser)))

	return api
}
