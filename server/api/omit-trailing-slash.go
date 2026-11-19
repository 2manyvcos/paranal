package api

import (
	"net/http"
	"net/url"
	"strings"
)

func OmitTrailingSlash(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		p := strings.TrimSuffix(req.URL.Path, "/")
		rp := strings.TrimSuffix(req.URL.RawPath, "/")
		if req.URL.Path != "/" && len(p) < len(req.URL.Path) && (req.URL.RawPath == "" || (len(rp) < len(req.URL.RawPath))) {
			req2 := new(http.Request)
			*req2 = *req
			req2.URL = new(url.URL)
			*req2.URL = *req.URL
			req2.URL.Path = p
			req2.URL.RawPath = rp
			req = req2
		}
		handler.ServeHTTP(res, req)
	})
}
