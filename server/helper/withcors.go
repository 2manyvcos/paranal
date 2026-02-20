package helper

import (
	"net/http"
	"strings"
)

func WithCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app := GetApp(req)

		if app.Config.CORSOrigin != "" {
			res.Header().Set("Access-Control-Allow-Origin", app.Config.CORSOrigin)
		}
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		allowedHeaders := []string{"Accept", "Content-Type", "Authorization"}
		if app.Config.Auth.RemoteUser.Enabled {
			allowedHeaders = append(allowedHeaders, app.Config.Auth.RemoteUser.HeaderName)
			if app.Config.Auth.RemoteUser.GroupsHeaderName != "" {
				allowedHeaders = append(allowedHeaders, app.Config.Auth.RemoteUser.GroupsHeaderName)
			}
		}
		res.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

		handler.ServeHTTP(res, req)
	})
}
