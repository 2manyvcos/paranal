package helper

import (
	"net/http"
	"strings"
)

func OmitTrailingSlash(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
		}

		handler.ServeHTTP(res, req)
	})
}
