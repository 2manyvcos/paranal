package api

import (
	"context"
	"net/http"
)

func WithApp(app *App, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), "app", app)))
	})
}

func GetApp(req *http.Request) *App {
	return req.Context().Value("app").(*App)
}
