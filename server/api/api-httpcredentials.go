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

func GetHTTPCredentials(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	records, err := app.ListHTTPCredentials()
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responsePayload := make([]struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Key      string `json:"key"`
		HasValue bool   `json:"hasValue"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].Name = record.Name
		responsePayload[i].Type = data.HTTP_CREDENTIAL_TYPE_NAMES[record.Type]
		responsePayload[i].Key = record.Key
		responsePayload[i].HasValue = record.Value != ""

	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(responsePayload)
}

func PostHTTPCredentials(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newRecord := data.HTTPCredential{
		Name: requestPayload.Name,
		Type: data.HTTP_CREDENTIAL_TYPE_CODES[requestPayload.Type],
		Key:  requestPayload.Key,
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
	err = app.CreateHTTPCredential(newRecord, false)
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

func GetHTTPCredentialsByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	name := req.PathValue("name")

	record, err := app.GetHTTPCredential(name)
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
		Name     string `json:"name"`
		Type     string `json:"type"`
		Key      string `json:"key"`
		HasValue bool   `json:"hasValue"`
	}{
		Name:     record.Name,
		Type:     data.HTTP_CREDENTIAL_TYPE_NAMES[record.Type],
		Key:      record.Key,
		HasValue: record.Value != "",
	})
}

func PatchHTTPCredentialsByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	name := req.PathValue("name")

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Type  utils.Optional[string] `json:"type"`
		Key   utils.Optional[string] `json:"key"`
		Value utils.Optional[string] `json:"value"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetHTTPCredential(name)
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
	if requestPayload.Type.IsDefined {
		updatedRecord.Type = data.HTTP_CREDENTIAL_TYPE_CODES[requestPayload.Type.Value]
	}
	requestPayload.Key.ApplyIfDefined(&updatedRecord.Key)
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
	err = app.UpdateHTTPCredential(updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func DeleteHTTPCredentialsByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	name := req.PathValue("name")

	err := app.DeleteHTTPCredential(name)
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
