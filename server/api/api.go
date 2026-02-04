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

	api.Handle("GET /users", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetUsers))))
	api.Handle("POST /users", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PostUsers))))
	api.Handle("GET /users/{userName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetUsersByName))))
	api.Handle("PATCH /users/{userName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PatchUsersByName))))
	api.Handle("DELETE /users/{userName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(DeleteUsersByName))))

	api.Handle("GET /httpcredentials", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetHTTPCredentials))))
	api.Handle("POST /httpcredentials", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PostHTTPCredentials))))
	api.Handle("GET /httpcredentials/{credentialName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetHTTPCredentialsByName))))
	api.Handle("PATCH /httpcredentials/{credentialName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PatchHTTPCredentialsByName))))
	api.Handle("DELETE /httpcredentials/{credentialName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(DeleteHTTPCredentialsByName))))

	api.Handle("GET /sshcredentials", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetSSHCredentials))))
	api.Handle("POST /sshcredentials", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PostSSHCredentials))))
	api.Handle("GET /sshcredentials/{credentialName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetSSHCredentialsByName))))
	api.Handle("PATCH /sshcredentials/{credentialName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PatchSSHCredentialsByName))))
	api.Handle("DELETE /sshcredentials/{credentialName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(DeleteSSHCredentialsByName))))

	api.Handle("GET /usercredentials", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetUserCredentials))))
	api.Handle("POST /usercredentials", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PostUserCredentials))))
	api.Handle("GET /usercredentials/{credentialName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetUserCredentialsByName))))
	api.Handle("PATCH /usercredentials/{credentialName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PatchUserCredentialsByName))))
	api.Handle("DELETE /usercredentials/{credentialName}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(DeleteUserCredentialsByName))))

	api.Handle("GET /services", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetServices))))
	api.Handle("POST /services", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PostServices))))
	api.Handle("GET /services/{serviceID}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetServicesByID))))
	api.Handle("PATCH /services/{serviceID}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PatchServicesByID))))
	api.Handle("DELETE /services/{serviceID}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(DeleteServicesByID))))
	api.Handle("POST /services/{serviceID}/favorite", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PostServicesByIDFavorite))))
	api.Handle("DELETE /services/{serviceID}/favorite", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(DeleteServicesByIDFavorite))))
	api.Handle("POST /services/{serviceID}/run-script", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PostServicesByIDRunScript))))
	api.Handle("GET /services/{serviceID}/scripts", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetServicesByIDScripts))))
	api.Handle("POST /services/{serviceID}/scripts", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PostServicesByIDScripts))))
	api.Handle("GET /services/{serviceID}/scripts/{scriptID}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(GetServicesByIDScriptsByID))))
	api.Handle("PATCH /services/{serviceID}/scripts/{scriptID}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PatchServicesByIDScriptsByID))))
	api.Handle("DELETE /services/{serviceID}/scripts/{scriptID}", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(DeleteServicesByIDScriptsByID))))

	api.Handle("POST /find-website-logo", helper.WithAuth(helper.WithRole(data.UserRoleAdmin, http.HandlerFunc(PostFindWebsiteLogo))))

	return api
}
