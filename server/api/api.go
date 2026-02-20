package api

import (
	"net/http"

	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/schema"
)

func New() http.Handler {
	api := http.NewServeMux()

	api.Handle("POST /v1/auth", http.HandlerFunc(PostAuth))

	api.Handle("GET /v1/events", helper.WithAuth(http.HandlerFunc(GetEvents)))

	api.Handle("GET /v1/user", helper.WithAuth(http.HandlerFunc(GetUser)))
	api.Handle("PATCH /v1/user", helper.WithAuth(http.HandlerFunc(PatchUser)))
	api.Handle("DELETE /v1/user", helper.WithAuth(http.HandlerFunc(DeleteUser)))
	api.Handle("PUT /v1/user/password", helper.WithAuth(http.HandlerFunc(PutUserPassword)))

	api.Handle("GET /v1/user/alertchannels", helper.WithAuth(http.HandlerFunc(GetUserAlertChannels)))
	api.Handle("POST /v1/user/alertchannels", helper.WithAuth(http.HandlerFunc(PostUserAlertChannels)))
	api.Handle("GET /v1/user/alertchannels/{channelID}", helper.WithAuth(http.HandlerFunc(GetUserAlertChannelsByID)))
	api.Handle("PATCH /v1/user/alertchannels/{channelID}", helper.WithAuth(http.HandlerFunc(PatchUserAlertChannelsByID)))
	api.Handle("DELETE /v1/user/alertchannels/{channelID}", helper.WithAuth(http.HandlerFunc(DeleteUserAlertChannelsByID)))

	api.Handle("GET /v1/layout", helper.WithAuth(http.HandlerFunc(GetLayout)))
	api.Handle("PUT /v1/layout", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PutLayout))))

	api.Handle("POST /v1/find-website-logo", helper.WithAuth(http.HandlerFunc(PostFindWebsiteLogo)))

	api.Handle("GET /v1/users", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetUsers))))
	api.Handle("POST /v1/users", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PostUsers))))
	api.Handle("GET /v1/users/{userName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetUsersByName))))
	api.Handle("PATCH /v1/users/{userName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PatchUsersByName))))
	api.Handle("DELETE /v1/users/{userName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(DeleteUsersByName))))

	api.Handle("GET /v1/httpcredentials", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetHTTPCredentials))))
	api.Handle("POST /v1/httpcredentials", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PostHTTPCredentials))))
	api.Handle("GET /v1/httpcredentials/{credentialName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetHTTPCredentialsByName))))
	api.Handle("PATCH /v1/httpcredentials/{credentialName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PatchHTTPCredentialsByName))))
	api.Handle("DELETE /v1/httpcredentials/{credentialName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(DeleteHTTPCredentialsByName))))

	api.Handle("GET /v1/sshcredentials", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetSSHCredentials))))
	api.Handle("POST /v1/sshcredentials", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PostSSHCredentials))))
	api.Handle("GET /v1/sshcredentials/{credentialName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetSSHCredentialsByName))))
	api.Handle("PATCH /v1/sshcredentials/{credentialName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PatchSSHCredentialsByName))))
	api.Handle("DELETE /v1/sshcredentials/{credentialName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(DeleteSSHCredentialsByName))))

	api.Handle("GET /v1/usercredentials", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetUserCredentials))))
	api.Handle("POST /v1/usercredentials", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PostUserCredentials))))
	api.Handle("GET /v1/usercredentials/{credentialName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetUserCredentialsByName))))
	api.Handle("PATCH /v1/usercredentials/{credentialName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PatchUserCredentialsByName))))
	api.Handle("DELETE /v1/usercredentials/{credentialName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(DeleteUserCredentialsByName))))

	api.Handle("GET /v1/services", helper.WithAuth(http.HandlerFunc(GetServices)))
	api.Handle("POST /v1/services", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PostServices))))
	api.Handle("GET /v1/services/{serviceID}", helper.WithAuth(http.HandlerFunc(GetServicesByID)))
	api.Handle("PATCH /v1/services/{serviceID}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PatchServicesByID))))
	api.Handle("DELETE /v1/services/{serviceID}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(DeleteServicesByID))))
	api.Handle("PATCH /v1/services/{serviceID}/config", helper.WithAuth(http.HandlerFunc(PatchServicesByIDConfig)))

	api.Handle("POST /v1/services/{serviceID}/run-script", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PostServicesByIDRunScript))))

	api.Handle("GET /v1/services/{serviceID}/scripts", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetServicesByIDScripts))))
	api.Handle("POST /v1/services/{serviceID}/scripts", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PostServicesByIDScripts))))
	api.Handle("GET /v1/services/{serviceID}/scripts/{scriptID}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetServicesByIDScriptsByID))))
	api.Handle("PATCH /v1/services/{serviceID}/scripts/{scriptID}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PatchServicesByIDScriptsByID))))
	api.Handle("DELETE /v1/services/{serviceID}/scripts/{scriptID}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(DeleteServicesByIDScriptsByID))))
	api.Handle("POST /v1/services/{serviceID}/scripts/{scriptID}/run", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PostServicesByIDScriptsByIDRun))))

	api.Handle("GET /v1/services/{serviceID}/actions", helper.WithAuth(http.HandlerFunc(GetServicesByIDActions)))
	api.Handle("POST /v1/services/{serviceID}/actions/{actionName}/run", helper.WithAuth(http.HandlerFunc(PostServicesByIDActionsByNameRun)))

	api.Handle("GET /v1/uptimestatuses", helper.WithAuth(http.HandlerFunc(GetUptimeStatuses)))
	api.Handle("GET /v1/services/{serviceID}/uptimestatuses", helper.WithAuth(http.HandlerFunc(GetServicesByIDUptimeStatuses)))

	api.Handle("GET /v1/versions", helper.WithAuth(http.HandlerFunc(GetVersions)))
	api.Handle("GET /v1/services/{serviceID}/versions", helper.WithAuth(http.HandlerFunc(GetServicesByIDVersions)))
	api.Handle("GET /v1/services/{serviceID}/versions/{versionName}/details", helper.WithAuth(http.HandlerFunc(GetServicesByIDVersionsByNameDetails)))

	api.Handle("GET /v1/maintenancetasks", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetMaintenanceTasks))))
	api.Handle("GET /v1/maintenancetasks/{taskName}", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(GetMaintenanceTasksByName))))
	api.Handle("POST /v1/maintenancetasks/{taskName}/run", helper.WithAuth(helper.WithRole(schema.UserRoleAdmin, http.HandlerFunc(PostMaintenanceTasksByNameRun))))

	return api
}
