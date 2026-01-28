package api

import (
	"net/http"

	"github.com/2manyvcos/paranal/server/data"
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
	api.Handle("POST /auth", http.HandlerFunc(PostAuth))
	api.Handle("GET /user", helper.WithAuth(http.HandlerFunc(GetUser)))
	api.Handle("GET /users", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetUsers))))
	api.Handle("GET /users/{username}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetUsersByUsername))))
	api.Handle("DELETE /users/{username}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(DeleteUsersByUsername))))

	return api
}
