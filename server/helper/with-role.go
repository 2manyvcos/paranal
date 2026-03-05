package helper

import (
	"net/http"
	"slices"
)

func WithRole(role int, handler http.Handler) http.Handler {
	return WithRoles([]int{role}, handler)
}

func WithRoles(roles []int, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authorizedUser := GetAuthorizedUser(req)

		if len(roles) > 0 && (authorizedUser == nil || !slices.Contains(roles, authorizedUser.Role)) {
			http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}

		handler.ServeHTTP(res, req)
	})
}
