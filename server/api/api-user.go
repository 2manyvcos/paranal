package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
)

func GetUser(res http.ResponseWriter, req *http.Request) {
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(struct {
		Name          string `json:"name"`
		Role          string `json:"role"`
		HasPassword   bool   `json:"hasPassword"`
		DisplayName   string `json:"displayName"`
		StartPage     string `json:"startPage"`
		ErrorAlerts   bool   `json:"errorAlerts"`
		UptimeAlerts  bool   `json:"uptimeAlerts"`
		VersionAlerts bool   `json:"versionAlerts"`
	}{
		Name:          authorizedUser.Name,
		Role:          schema.UserRoleNames[authorizedUser.Role],
		HasPassword:   authorizedUser.PasswordHash != "",
		DisplayName:   authorizedUser.DisplayName,
		StartPage:     authorizedUser.StartPage,
		ErrorAlerts:   authorizedUser.ErrorAlerts,
		UptimeAlerts:  authorizedUser.UptimeAlerts,
		VersionAlerts: authorizedUser.VersionAlerts,
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PatchUser(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		DisplayName   utils.Optional[string] `json:"displayName"`
		StartPage     utils.Optional[string] `json:"startPage"`
		ErrorAlerts   utils.Optional[bool]   `json:"errorAlerts"`
		UptimeAlerts  utils.Optional[bool]   `json:"uptimeAlerts"`
		VersionAlerts utils.Optional[bool]   `json:"versionAlerts"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	updatedRecord := *authorizedUser
	requestPayload.DisplayName.ApplyIfDefined(&updatedRecord.DisplayName)
	requestPayload.StartPage.ApplyIfDefined(&updatedRecord.StartPage)
	requestPayload.ErrorAlerts.ApplyIfDefined(&updatedRecord.ErrorAlerts)
	requestPayload.UptimeAlerts.ApplyIfDefined(&updatedRecord.UptimeAlerts)
	requestPayload.VersionAlerts.ApplyIfDefined(&updatedRecord.VersionAlerts)
	err = app.UpdateUsers(schema.UserQuery{Name: &authorizedUser.Name}, updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(authorizedUser.Name, "/v1/user"))
}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err := app.DeleteUsers(schema.UserQuery{Name: &authorizedUser.Name})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error deleting record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(authorizedUser.Name, "/v1/user"))
	app.PublishClientEvent(schema.NewUpdateEvent("/v1/users/" + url.PathEscape(authorizedUser.Name)))
}

func PutUserPassword(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if requestPayload.NewPassword == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if authorizedUser.PasswordHash != "" {
		matches, err := crypto.CompareToHash(requestPayload.CurrentPassword, authorizedUser.PasswordHash)
		if !matches {
			if err != nil {
				log.Printf("Error comparing hash - %s\n", err)
			}
			maskAuthRejection()
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	updatedRecord := *authorizedUser
	updatedRecord.PasswordHash, err = crypto.Hash(requestPayload.NewPassword)
	if err != nil {
		log.Printf("Error generating hash - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := updatedRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.UpdateUsers(schema.UserQuery{Name: &authorizedUser.Name}, updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(authorizedUser.Name, "/v1/user"))
	app.PublishClientEvent(schema.NewUpdateEvent("/v1/users/" + url.PathEscape(authorizedUser.Name)))
}
