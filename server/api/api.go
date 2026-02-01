package api

import (
	"net/http"

	"github.com/2manyvcos/paranal/server/data"
	"github.com/2manyvcos/paranal/server/helper"
)

const API_VERSION = "v1"

func New() http.Handler {
	api := http.NewServeMux()

	api.Handle("POST /auth", http.HandlerFunc(PostAuth))

	api.Handle("GET /user", helper.WithAuth(http.HandlerFunc(GetUser)))
	api.Handle("PATCH /user", helper.WithAuth(http.HandlerFunc(PatchUser)))
	api.Handle("PUT /user/password", helper.WithAuth(http.HandlerFunc(PutUserPassword)))

	api.Handle("GET /users", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetUsers))))
	api.Handle("POST /users", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(PostUsers))))
	api.Handle("GET /users/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetUsersByName))))
	api.Handle("PATCH /users/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(PatchUsersByName))))
	api.Handle("DELETE /users/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(DeleteUsersByName))))

	api.Handle("GET /httpcredentials", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetHTTPCredentials))))
	api.Handle("POST /httpcredentials", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(PostHTTPCredentials))))
	api.Handle("GET /httpcredentials/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetHTTPCredentialsByName))))
	api.Handle("PATCH /httpcredentials/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(PatchHTTPCredentialsByName))))
	api.Handle("DELETE /httpcredentials/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(DeleteHTTPCredentialsByName))))

	api.Handle("GET /sshcredentials", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetSSHCredentials))))
	api.Handle("POST /sshcredentials", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(PostSSHCredentials))))
	api.Handle("GET /sshcredentials/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetSSHCredentialsByName))))
	api.Handle("PATCH /sshcredentials/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(PatchSSHCredentialsByName))))
	api.Handle("DELETE /sshcredentials/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(DeleteSSHCredentialsByName))))

	api.Handle("GET /usercredentials", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetUserCredentials))))
	api.Handle("POST /usercredentials", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(PostUserCredentials))))
	api.Handle("GET /usercredentials/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(GetUserCredentialsByName))))
	api.Handle("PATCH /usercredentials/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(PatchUserCredentialsByName))))
	api.Handle("DELETE /usercredentials/{name}", helper.WithAuth(helper.WithRole(data.USER_ROLE_ADMIN, http.HandlerFunc(DeleteUserCredentialsByName))))

	return api
}
