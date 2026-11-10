package api

import (
	"net/http"
)

func New() http.Handler {
	api := http.NewServeMux()

	api.Handle("GET /v1/events", http.HandlerFunc(GetEvents))

	return api
}
