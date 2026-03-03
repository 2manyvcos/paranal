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

		if req.Method == http.MethodOptions {
			allowedHeaders := []string{"Accept", "Content-Type", "Authorization"}
			if app.Config.Auth.RemoteUser.Enabled {
				allowedHeaders = append(allowedHeaders, app.Config.Auth.RemoteUser.HeaderName)
				if app.Config.Auth.RemoteUser.GroupsHeaderName != "" {
					allowedHeaders = append(allowedHeaders, app.Config.Auth.RemoteUser.GroupsHeaderName)
				}
			}
			res.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

			corsRes := corsResponseWriter{}
			corsReq := req.Clone(req.Context())
			handler.ServeHTTP(&corsRes, corsReq)
			if allow := corsRes.header.Get("Allow"); allow != "" {
				res.Header().Set("Access-Control-Allow-Methods", allow)
			}

			res.WriteHeader(http.StatusNoContent)
			return
		}

		handler.ServeHTTP(res, req)
	})
}

type corsResponseWriter struct {
	header http.Header
}

func (res *corsResponseWriter) Header() http.Header {
	if res.header == nil {
		res.header = make(http.Header)
	}
	return res.header
}

func (res *corsResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (res *corsResponseWriter) WriteHeader(int) {}
