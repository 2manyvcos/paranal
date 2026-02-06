package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/utils"
)

func GetUser(res http.ResponseWriter, req *http.Request) {
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Role        string `json:"role"`
		HasPassword bool   `json:"hasPassword"`
	}{
		Name:        authorizedUser.Name,
		DisplayName: authorizedUser.DisplayName,
		Role:        schema.UserRoleNames[authorizedUser.Role],
		HasPassword: authorizedUser.PasswordHash != "",
	})
}

func PatchUser(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil || authorizedUser.Name == "" {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		DisplayName utils.Optional[string] `json:"displayName"`
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
	err = app.UpdateUsers(schema.UserQuery{Name: &authorizedUser.Name}, updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func PutUserPassword(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil || authorizedUser.Name == "" {
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
}
