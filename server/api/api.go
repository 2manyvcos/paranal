package api

import (
	"net/http"
)

func New() http.Handler {
	api := http.NewServeMux()

	return api
}
