package helper

import (
	"context"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/2manyvcos/paranal/utils"
)

func WithAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authorizedUser := authorize(req)

		if authorizedUser == nil {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), "authorizedUser", authorizedUser)))
	})
}

func GetAuthorizedUser(req *http.Request) *schema.User {
	return req.Context().Value("authorizedUser").(*schema.User)
}

func authorize(req *http.Request) *schema.User {
	if authorizedUser := authorizeRemoteUser(req); authorizedUser != nil {
		return authorizedUser
	}

	return authorizeBearer(req)
}

func authorizeRemoteUser(req *http.Request) *schema.User {
	app := GetApp(req)

	if !app.Config.Auth.RemoteUser.Enabled {
		return nil
	}

	userName := req.Header.Get(app.Config.Auth.RemoteUser.HeaderName)
	if userName == "" {
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

	var authorizedUser *schema.User
	if user, err := app.GetUser(schema.UserQuery{Name: &userName}); err != nil {
		if !errors.Is(err, schema.ErrNotFound) {
			log.Printf("Error loading user - %s\n", err)
			return nil
		}
	} else {
		authorizedUser = &user
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
		newUser := schema.User{
			Name: userName,
			Role: map[bool]int{false: schema.UserRoleCommon, true: schema.UserRoleAdmin}[admin],
		}
		err := app.CreateUser(newUser)
		if err != nil {
			log.Printf("Error creating user - %s\n", err)
			return nil
		}
		return &newUser
	}

	if admin && authorizedUser.Role != schema.UserRoleAdmin {
		newUser := *authorizedUser
		newUser.Role = schema.UserRoleAdmin
		err := app.UpdateUsers(schema.UserQuery{Name: &userName}, newUser)
		if err != nil {
			log.Printf("Error updating user - %s\n", err)
		} else {
			authorizedUser = &newUser
		}
	}

	return authorizedUser
}

func authorizeBearer(req *http.Request) *schema.User {
	app := GetApp(req)

	bearer := req.Header.Get("Authorization")
	if !strings.HasPrefix(bearer, "Bearer ") {
		return nil
	}

	ok, userName, err := crypto.ValidateJWTToken(app.Config.Auth.JWT.Secret, strings.TrimSpace(strings.TrimPrefix(bearer, "Bearer ")))
	if !ok || userName == "" || err != nil {
		// if err != nil {
		// 	log.Printf("Error validating token - %s\n", err)
		// }
		return nil
	}

	user, err := app.GetUser(schema.UserQuery{Name: &userName})
	if errors.Is(err, schema.ErrNotFound) {
		return nil
	}
	if err != nil {
		log.Printf("Error loading user - %s\n", err)
		return nil
	}
	return &user
}
