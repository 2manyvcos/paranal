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

func GetSSHCredentials(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	records, err := app.ListSSHCredentials()
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responsePayload := make([]struct {
		Name          string `json:"name"`
		User          string `json:"user"`
		HasPassword   bool   `json:"hasPassword"`
		HasPrivateKey bool   `json:"hasPrivateKey"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].Name = record.Name
		responsePayload[i].User = record.User
		responsePayload[i].HasPassword = record.Password != ""
		responsePayload[i].HasPrivateKey = record.PrivateKey != ""

	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(responsePayload)
}

func PostSSHCredentials(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Name       string `json:"name"`
		User       string `json:"user"`
		Password   string `json:"password"`
		PrivateKey string `json:"privateKey"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newRecord := data.SSHCredential{
		Name: requestPayload.Name,
		User: requestPayload.User,
	}
	if requestPayload.Password != "" {
		newRecord.Password, err = crypto.Encrypt(app.Config.SecretKey, requestPayload.Password)
		if err != nil {
			log.Printf("Error encrypting value - %s\n", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if requestPayload.PrivateKey != "" {
		newRecord.PrivateKey, err = crypto.Encrypt(app.Config.SecretKey, requestPayload.PrivateKey)
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
	err = app.CreateSSHCredential(newRecord, false)
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

func GetSSHCredentialsByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	name := req.PathValue("name")

	record, err := app.GetSSHCredential(name)
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
		Name          string `json:"name"`
		User          string `json:"user"`
		HasPassword   bool   `json:"hasPassword"`
		HasPrivateKey bool   `json:"hasPrivateKey"`
	}{
		Name:          record.Name,
		User:          record.User,
		HasPassword:   record.Password != "",
		HasPrivateKey: record.PrivateKey != "",
	})
}

func PatchSSHCredentialsByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	name := req.PathValue("name")

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		User       utils.Optional[string] `json:"user"`
		Password   utils.Optional[string] `json:"password"`
		PrivateKey utils.Optional[string] `json:"privateKey"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetSSHCredential(name)
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
	requestPayload.User.ApplyIfDefined(&updatedRecord.User)
	if requestPayload.Password.IsDefined {
		if requestPayload.Password.Value == "" {
			updatedRecord.Password = ""
		} else {
			updatedRecord.Password, err = crypto.Encrypt(app.Config.SecretKey, requestPayload.Password.Value)
			if err != nil {
				log.Printf("Error encrypting value - %s\n", err)
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}
	if requestPayload.PrivateKey.IsDefined {
		if requestPayload.PrivateKey.Value == "" {
			updatedRecord.PrivateKey = ""
		} else {
			updatedRecord.PrivateKey, err = crypto.Encrypt(app.Config.SecretKey, requestPayload.PrivateKey.Value)
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
	err = app.UpdateSSHCredential(updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func DeleteSSHCredentialsByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	name := req.PathValue("name")

	err := app.DeleteSSHCredential(name)
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
