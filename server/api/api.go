package api

import (
	"net/http"
)

func New() http.Handler {
	api := http.NewServeMux()

	api.Handle("GET /v1/events", http.HandlerFunc(GetEvents))

	api.Handle("GET /v1/registries/{registryName}/states", http.HandlerFunc(GetRegistriesByNameStates))
	api.Handle("GET /v1/registries/{registryName}/states/{stateID}", http.HandlerFunc(GetRegistriesByNameStatesByID))

	return api
}
