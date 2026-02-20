package helper

import (
	"net/http"
)

func WithCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app := GetApp(req)

		if app.Config.AppURL != "" {
			res.Header().Set("Access-Control-Allow-Origin", app.Config.AppURL)
		}
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		res.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")

		handler.ServeHTTP(res, req)
	})
}
