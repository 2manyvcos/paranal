package helper

import (
	"context"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/2manyvcos/paranal/utils"
)

func WithAuth(handler http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		authorizedUser := authorize(req)

		if authorizedUser == nil {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), "authorizedUser", authorizedUser)))
	}
}

func GetAuthorizedUser(req *http.Request) *data.User {
	return req.Context().Value("authorizedUser").(*data.User)
}

func authorize(req *http.Request) *data.User {
	if authorizedUser := authorizeRemoteUser(req); authorizedUser != nil {
		return authorizedUser
	}

	return authorizeBearer(req)
}

func authorizeRemoteUser(req *http.Request) *data.User {
	app := GetApp(req)

	if !app.Config.Auth.RemoteUser.Enabled {
		return nil
	}

	username := req.Header.Get(app.Config.Auth.RemoteUser.HeaderName)
	if username == "" {
		return nil
	}

	clientIP := req.RemoteAddr
	if portSeparator := strings.LastIndex(clientIP, ":"); portSeparator >= 0 && strings.LastIndex(clientIP, "]") < portSeparator {
		clientIP = clientIP[:portSeparator]
	}
	if len(app.Config.Auth.RemoteUser.Whitelist) > 0 && !slices.Contains(app.Config.Auth.RemoteUser.Whitelist, clientIP) {
		log.Printf("Remote-User auth attempt blocked for untrusted client \"%s\"\n", clientIP)
		return nil
	}

	authorizedUser, err := app.GetUser(username)
	if err != nil {
		log.Printf("Error loading user - %s\n", err)
		return nil
	}

	var admin bool
	if app.Config.Auth.RemoteUser.GroupsHeaderName != "" && app.Config.Auth.RemoteUser.AdminGroup != "" {
		groups := utils.ParseList(req.Header.Get(app.Config.Auth.RemoteUser.GroupsHeaderName), nil)
		admin = slices.Contains(groups, app.Config.Auth.RemoteUser.AdminGroup)
	}

	if authorizedUser == nil {
		if !app.Config.Auth.RemoteUser.CreateUnknownUsers {
			return nil
		}
		role := data.USER_ROLE_COMMON
		if admin {
			role = data.USER_ROLE_ADMIN
		}
		newUser := data.User{Name: username, Role: role}
		err := app.InsertUser(newUser)
		if err != nil {
			log.Printf("Error creating user - %s\n", err)
			return nil
		}
		return &newUser
	}

	if admin && authorizedUser.Role != data.USER_ROLE_ADMIN {
		newUser := *authorizedUser
		newUser.Role = data.USER_ROLE_ADMIN
		err := app.UpdateUser(newUser)
		if err != nil {
			log.Printf("Error updating user - %s\n", err)
		} else {
			authorizedUser = &newUser
		}
	}

	return authorizedUser
}

func authorizeBearer(req *http.Request) *data.User {
	app := GetApp(req)

	bearer := req.Header.Get("Authorization")
	if !strings.HasPrefix(bearer, "Bearer ") {
		return nil
	}

	ok, username, err := crypto.JWTValidateToken(app.Config.Auth.JWT.Secret, strings.TrimSpace(strings.TrimPrefix(bearer, "Bearer ")))
	if !ok || err != nil {
		// if err != nil {
		// 	log.Printf("Error validating token - %s\n", err)
		// }
		return nil
	}

	user, err := app.GetUser(username)
	if err != nil {
		log.Printf("Error loading user - %s\n", err)
		return nil
	}
	return user
}
