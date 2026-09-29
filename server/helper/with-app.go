package helper

import (
	"context"
	"net/http"

	"github.com/2manyvcos/paranal/server/application"
)

func WithApp(app *application.App, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), "app", app)))
	})
}

func GetApp(req *http.Request) *application.App {
	return req.Context().Value("app").(*application.App)
}
