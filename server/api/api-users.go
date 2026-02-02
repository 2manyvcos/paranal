package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/utils"
)

func GetUsers(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	records, err := app.ListUsers()
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responsePayload := make([]struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Role        string `json:"role"`
		HasPassword bool   `json:"hasPassword"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].Name = record.Name
		responsePayload[i].DisplayName = record.DisplayName
		responsePayload[i].Role = data.USER_ROLE_NAMES[record.Role]
		responsePayload[i].HasPassword = record.PasswordHash != ""
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(responsePayload)
}

func PostUsers(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Role        string `json:"role"`
		Password    string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newRecord := data.User{
		Name:        requestPayload.Name,
		DisplayName: requestPayload.DisplayName,
		Role:        data.USER_ROLE_CODES[requestPayload.Role],
	}
	if requestPayload.Password != "" {
		newRecord.PasswordHash, err = crypto.Hash(requestPayload.Password)
		if err != nil {
			log.Printf("Error generating hash - %s\n", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if err := newRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.CreateUser(newRecord, false)
	if errors.Is(err, data.ErrConflict) {
		http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Error creating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
}

func GetUsersByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	name := req.PathValue("name")

	record, err := app.GetUser(name)
	if errors.Is(err, data.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Role        string `json:"role"`
		HasPassword bool   `json:"hasPassword"`
	}{
		Name:        record.Name,
		DisplayName: record.DisplayName,
		Role:        data.USER_ROLE_NAMES[record.Role],
		HasPassword: record.PasswordHash != "",
	})
}

func PatchUsersByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	name := req.PathValue("name")

	if authorizedUser != nil && authorizedUser.Name == name {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		DisplayName utils.Optional[string] `json:"displayName"`
		Role        utils.Optional[string] `json:"role"`
		Password    utils.Optional[string] `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUser(name)
	if errors.Is(err, data.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	updatedRecord := record
	requestPayload.DisplayName.ApplyIfDefined(&updatedRecord.DisplayName)
	if requestPayload.Role.IsDefined {
		updatedRecord.Role = data.USER_ROLE_CODES[requestPayload.Role.Value]
	}
	if requestPayload.Password.IsDefined {
		if requestPayload.Password.Value == "" {
			updatedRecord.PasswordHash = ""
		} else {
			updatedRecord.PasswordHash, err = crypto.Hash(requestPayload.Password.Value)
			if err != nil {
				log.Printf("Error generating hash - %s\n", err)
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}
	if err := updatedRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.UpdateUser(updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func DeleteUsersByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	name := req.PathValue("name")

	if authorizedUser != nil && authorizedUser.Name == name {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err := app.DeleteUser(name)
	if errors.Is(err, data.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error deleting record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
