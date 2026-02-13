package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
)

type Layout struct {
	Pages []struct {
		Name   string `json:"name"`
		Header string `json:"header"`

		Sections []struct {
			Name       string   `json:"name"`
			Icon       string   `json:"icon"`
			ServiceIDs []string `json:"serviceIDs,omitempty"`
		} `json:"sections,omitempty"`
	} `json:"pages,omitempty"`
}

func GetLayout(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	setting, err := app.GetSetting("layout")
	if err != nil && !errors.Is(err, schema.ErrNotFound) {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var responsePayload Layout
	if setting != "" {
		err = json.Unmarshal([]byte(setting), &responsePayload)
		if err != nil {
			log.Printf("Error parsing layout - %s\n", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PutLayout(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload Layout
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	setting, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error encoding layout - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = app.CreateOrUpdateSetting("layout", string(setting))
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
