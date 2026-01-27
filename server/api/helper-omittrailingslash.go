package api

import (
	"net/http"
	"strings"
)

func OmitTrailingSlash(handler http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
		}

		handler.ServeHTTP(res, req)
	}
}
