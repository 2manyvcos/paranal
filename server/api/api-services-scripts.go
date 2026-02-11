package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/utils"
	"github.com/google/uuid"
)

func GetServicesByIDScripts(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	query := schema.ServiceScriptQuery{ServiceID: &serviceID}
	q := req.URL.Query()
	if v, ok := utils.LoadQueryValue(q, "id"); ok {
		query.ID = &v
	}
	if v, ok := utils.LoadQueryValue(q, "name"); ok {
		query.Name = &v
	}
	if v, ok := utils.LoadQueryValue(q, "schedule"); ok {
		query.Schedule = &v
	}
	if v, ok := utils.LoadQueryValue(q, "source"); ok {
		query.Source = &v
	}

	records, err := app.ListServiceScripts(&query)
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	scriptStates := app.State.GetServiceScriptStates(serviceID)

	responsePayload := make([]struct {
		ID       string     `json:"id"`
		Name     string     `json:"name"`
		Schedule string     `json:"schedule"`
		Source   string     `json:"source"`
		LastRun  *time.Time `json:"lastRun"`
		Error    *string    `json:"error"`
		NextRun  *time.Time `json:"nextRun"`
	}, len(records))
	for i, record := range records {
		state := scriptStates[record.ID]
		responsePayload[i].ID = record.ID
		responsePayload[i].Name = record.Name
		responsePayload[i].Schedule = record.Schedule
		responsePayload[i].Source = record.Source
		responsePayload[i].LastRun = state.LastRun
		if state.Error != nil {
			errorMessage := state.Error.Error()
			responsePayload[i].Error = &errorMessage
		}
		responsePayload[i].NextRun = state.NextRun
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PostServicesByIDScripts(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

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

	_, err = app.GetService(schema.ServiceQuery{ID: &serviceID})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newRecord := schema.ServiceScript{
		ID:        uuid.NewString(),
		Name:      requestPayload.Name,
		Schedule:  requestPayload.Schedule,
		Source:    requestPayload.Source,
		ServiceID: serviceID,
	}
	if err := newRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.CreateServiceScript(newRecord)
	if errors.Is(err, schema.ErrConflict) {
		http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Error creating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.State.OnServiceScriptChanged(newRecord)

	res.WriteHeader(http.StatusCreated)
}

func GetServicesByIDScriptsByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	scriptID := req.PathValue("scriptID")
	if serviceID == "" || scriptID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetServiceScript(schema.ServiceScriptQuery{ID: &scriptID, ServiceID: &serviceID})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	state := app.State.GetServiceScriptState(serviceID, scriptID)
	if state == nil {
		state = new(application.ServiceScriptState)
	}

	var error *string
	if state.Error != nil {
		errorMessage := state.Error.Error()
		error = &errorMessage
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(struct {
		ID       string     `json:"id"`
		Name     string     `json:"name"`
		Schedule string     `json:"schedule"`
		Source   string     `json:"source"`
		LastRun  *time.Time `json:"lastRun"`
		Error    *string    `json:"error"`
		NextRun  *time.Time `json:"nextRun"`
	}{
		ID:       record.ID,
		Name:     record.Name,
		Schedule: record.Schedule,
		Source:   record.Source,
		LastRun:  state.LastRun,
		Error:    error,
		NextRun:  state.NextRun,
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PatchServicesByIDScriptsByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	scriptID := req.PathValue("scriptID")
	if serviceID == "" || scriptID == "" {
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
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetServiceScript(schema.ServiceScriptQuery{ID: &scriptID, ServiceID: &serviceID})
	if errors.Is(err, schema.ErrNotFound) {
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
	err = app.UpdateServiceScripts(schema.ServiceScriptQuery{ID: &scriptID, ServiceID: &serviceID}, updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.State.OnServiceScriptChanged(updatedRecord)
}

func DeleteServicesByIDScriptsByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	scriptID := req.PathValue("scriptID")
	if serviceID == "" || scriptID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := app.DeleteServiceScripts(schema.ServiceScriptQuery{ID: &scriptID, ServiceID: &serviceID})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error deleting record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.State.OnServiceScriptDeleted(serviceID, scriptID)
}

func PostServicesByIDScriptsByIDRun(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	scriptID := req.PathValue("scriptID")
	if serviceID == "" || scriptID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ok := app.State.RunServiceScript(serviceID, scriptID)
	if !ok {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}
