package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/2manyvcos/paranal/server/data"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/scripts"
	"github.com/2manyvcos/paranal/utils"
	"github.com/google/uuid"
)

func GetServices(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	records, err := app.ListServices()
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responsePayload := make([]struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
		URL         string `json:"url"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].ID = record.ID
		responsePayload[i].Name = record.Name
		responsePayload[i].Description = record.Description
		responsePayload[i].Logo = record.Logo
		responsePayload[i].URL = record.URL
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(responsePayload)
}

func PostServices(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
		URL         string `json:"url"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newRecord := data.Service{
		ID:          uuid.New().String(),
		Name:        requestPayload.Name,
		Description: requestPayload.Description,
		Logo:        requestPayload.Logo,
		URL:         requestPayload.URL,
	}
	if err := newRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.CreateService(newRecord, false)
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

func GetServicesByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")

	record, err := app.GetService(serviceID)
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
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
		URL         string `json:"url"`
	}{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
		Logo:        record.Logo,
		URL:         record.URL,
	})
}

func PatchServicesByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Name        utils.Optional[string] `json:"name"`
		Description utils.Optional[string] `json:"description"`
		Logo        utils.Optional[string] `json:"logo"`
		URL         utils.Optional[string] `json:"url"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetService(serviceID)
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
	requestPayload.Name.ApplyIfDefined(&updatedRecord.Name)
	requestPayload.Description.ApplyIfDefined(&updatedRecord.Description)
	requestPayload.Logo.ApplyIfDefined(&updatedRecord.Logo)
	requestPayload.URL.ApplyIfDefined(&updatedRecord.URL)
	if err := updatedRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.UpdateService(updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func DeleteServicesByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")

	err := app.DeleteService(serviceID)
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

func PostServicesByIDRunScript(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Source string `json:"source"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = app.GetService(serviceID)
	if errors.Is(err, data.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	results, err := scripts.Run(app, serviceID, requestPayload.Source)
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(struct {
		Success bool  `json:"success"`
		Results []any `json:"results"`
	}{
		Success: true,
		Results: results,
	})
}

func GetServicesByIDScripts(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")

	records, err := app.ListScriptsByService(serviceID)
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responsePayload := make([]struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Schedule string `json:"schedule"`
		Source   string `json:"source"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].ID = strconv.Itoa(record.ID)
		responsePayload[i].Name = record.Name
		responsePayload[i].Schedule = record.Schedule
		responsePayload[i].Source = record.Source
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(responsePayload)
}

func PostServicesByIDScripts(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Name     string `json:"name"`
		Schedule string `json:"schedule"`
		Source   string `json:"source"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = app.GetService(serviceID)
	if errors.Is(err, data.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newRecord := data.Script{
		Name:      requestPayload.Name,
		Schedule:  requestPayload.Schedule,
		Source:    requestPayload.Source,
		ServiceID: serviceID,
	}
	if err := newRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.CreateScript(newRecord, false)
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

func GetServicesByIDScriptsByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	scriptID, err := strconv.Atoi(req.PathValue("scriptID"))
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetScriptByService(serviceID, scriptID)
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
		ID       string `json:"id"`
		Name     string `json:"name"`
		Schedule string `json:"schedule"`
	}{
		ID:       strconv.Itoa(record.ID),
		Name:     record.Name,
		Schedule: record.Schedule,
	})
}

func PatchServicesByIDScriptsByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	scriptID, err := strconv.Atoi(req.PathValue("scriptID"))
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Name     utils.Optional[string] `json:"name"`
		Schedule utils.Optional[string] `json:"schedule"`
		Source   utils.Optional[string] `json:"source"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetScriptByService(serviceID, scriptID)
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
	requestPayload.Name.ApplyIfDefined(&updatedRecord.Name)
	requestPayload.Schedule.ApplyIfDefined(&updatedRecord.Schedule)
	requestPayload.Source.ApplyIfDefined(&updatedRecord.Source)
	if err := updatedRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.UpdateScriptByService(updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func DeleteServicesByIDScriptsByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	scriptID, err := strconv.Atoi(req.PathValue("scriptID"))
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = app.DeleteScriptByService(serviceID, scriptID)
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
