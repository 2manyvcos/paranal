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

func GetUserCredentials(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	records, err := app.ListUserCredentials()
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responsePayload := make([]struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		HasValue    bool   `json:"hasValue"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].Name = record.Name
		responsePayload[i].Description = record.Description
		responsePayload[i].HasValue = record.Value != ""
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(responsePayload)
}

func PostUserCredentials(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Value       string `json:"value"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newRecord := data.UserCredential{
		Name:        requestPayload.Name,
		Description: requestPayload.Description,
	}
	if requestPayload.Value != "" {
		newRecord.Value, err = crypto.Encrypt(app.Config.SecretKey, requestPayload.Value)
		if err != nil {
			log.Printf("Error encrypting value - %s\n", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if err := newRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.CreateUserCredential(newRecord, false)
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

func GetUserCredentialsByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	credentialName := req.PathValue("credentialName")

	record, err := app.GetUserCredential(credentialName)
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
		Description string `json:"description"`
		HasValue    bool   `json:"hasValue"`
	}{
		Name:        record.Name,
		Description: record.Description,
		HasValue:    record.Value != "",
	})
}

func PatchUserCredentialsByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	credentialName := req.PathValue("credentialName")

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Description utils.Optional[string] `json:"description"`
		Value       utils.Optional[string] `json:"value"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUserCredential(credentialName)
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
	requestPayload.Description.ApplyIfDefined(&updatedRecord.Description)
	if requestPayload.Value.IsDefined {
		if requestPayload.Value.Value == "" {
			updatedRecord.Value = ""
		} else {
			updatedRecord.Value, err = crypto.Encrypt(app.Config.SecretKey, requestPayload.Value.Value)
			if err != nil {
				log.Printf("Error encrypting value - %s\n", err)
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}
	if err := updatedRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.UpdateUserCredential(updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func DeleteUserCredentialsByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	credentialName := req.PathValue("credentialName")

	err := app.DeleteUserCredential(credentialName)
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
