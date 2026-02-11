package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/scripts"
	"github.com/2manyvcos/paranal/utils"
	"github.com/google/uuid"
)

func GetServices(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil || authorizedUser.Name == "" {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	query := new(schema.UserServiceQuery)
	q := req.URL.Query()
	if q.Has("hidden") {
		v := q.Get("hidden")
		if v == "" {
			v = "true"
		}
		if p, err := strconv.ParseBool(v); err == nil {
			query.Hidden = &p
		} else {
			log.Printf("Error parsing query parameter - %s\n", err)
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	records, err := app.ListUserServices(authorizedUser.Name, query)
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
		Config      struct {
			Favorite      bool `json:"favorite"`
			Hidden        bool `json:"hidden"`
			UptimeAlerts  bool `json:"uptimeAlerts"`
			VersionAlerts bool `json:"versionAlerts"`
		} `json:"config"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].ID = record.ID
		responsePayload[i].Name = record.Name
		responsePayload[i].Description = record.Description
		responsePayload[i].Logo = record.Logo
		responsePayload[i].URL = record.URL
		responsePayload[i].Config.Favorite = record.Favorite
		responsePayload[i].Config.Hidden = record.Hidden
		responsePayload[i].Config.UptimeAlerts = record.UptimeAlerts
		responsePayload[i].Config.VersionAlerts = record.VersionAlerts
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
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

	newRecord := schema.Service{
		ID:          uuid.NewString(),
		Name:        requestPayload.Name,
		Description: requestPayload.Description,
		Logo:        requestPayload.Logo,
		URL:         requestPayload.URL,
	}
	if err := newRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.CreateService(newRecord)
	if errors.Is(err, schema.ErrConflict) {
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
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil || authorizedUser.Name == "" {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUserService(authorizedUser.Name, schema.UserServiceQuery{ID: &serviceID})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
		URL         string `json:"url"`
		Config      struct {
			Favorite      bool `json:"favorite"`
			Hidden        bool `json:"hidden"`
			UptimeAlerts  bool `json:"uptimeAlerts"`
			VersionAlerts bool `json:"versionAlerts"`
		} `json:"config"`
	}{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
		Logo:        record.Logo,
		URL:         record.URL,
		Config: struct {
			Favorite      bool `json:"favorite"`
			Hidden        bool `json:"hidden"`
			UptimeAlerts  bool `json:"uptimeAlerts"`
			VersionAlerts bool `json:"versionAlerts"`
		}{
			Favorite:      record.Favorite,
			Hidden:        record.Hidden,
			UptimeAlerts:  record.UptimeAlerts,
			VersionAlerts: record.VersionAlerts,
		},
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PatchServicesByID(res http.ResponseWriter, req *http.Request) {
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

	record, err := app.GetService(schema.ServiceQuery{ID: &serviceID})
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
	requestPayload.Description.ApplyIfDefined(&updatedRecord.Description)
	requestPayload.Logo.ApplyIfDefined(&updatedRecord.Logo)
	requestPayload.URL.ApplyIfDefined(&updatedRecord.URL)
	if err := updatedRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.UpdateServices(schema.ServiceQuery{ID: &serviceID}, updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func DeleteServicesByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := app.DeleteServices(schema.ServiceQuery{ID: &serviceID})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error deleting record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.State.OnServiceDeleted(serviceID)
}

func PatchServicesByIDConfig(res http.ResponseWriter, req *http.Request) {
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

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUserService(authorizedUser.Name, schema.UserServiceQuery{ID: &serviceID})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var requestPayload struct {
		Favorite      utils.Optional[bool] `json:"favorite"`
		Hidden        utils.Optional[bool] `json:"hidden"`
		UptimeAlerts  utils.Optional[bool] `json:"uptimeAlerts"`
		VersionAlerts utils.Optional[bool] `json:"versionAlerts"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	updatedRecord := schema.ServiceUserConfig{
		UserName:      authorizedUser.Name,
		ServiceID:     serviceID,
		Favorite:      record.Favorite,
		Hidden:        record.Hidden,
		UptimeAlerts:  record.UptimeAlerts,
		VersionAlerts: record.VersionAlerts,
	}
	requestPayload.Favorite.ApplyIfDefined(&updatedRecord.Favorite)
	requestPayload.Hidden.ApplyIfDefined(&updatedRecord.Hidden)
	requestPayload.UptimeAlerts.ApplyIfDefined(&updatedRecord.UptimeAlerts)
	requestPayload.VersionAlerts.ApplyIfDefined(&updatedRecord.VersionAlerts)
	if err := updatedRecord.Valid(); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.CreateOrUpdateServiceUserConfig(updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func PostServicesByIDRunScript(res http.ResponseWriter, req *http.Request) {
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
		Source string `json:"source"`
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

	results, err := scripts.RunServiceScript(app, serviceID, requestPayload.Source)
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(res).Encode(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}{
			Success: false,
			Error:   err.Error(),
		})
		if err != nil {
			log.Printf("Error encoding response payload - %s\n", err)
		}
		return
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(struct {
		Success bool  `json:"success"`
		Results []any `json:"results"`
	}{
		Success: true,
		Results: results,
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func GetServicesByIDScripts(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	records, err := app.ListServiceScripts(&schema.ServiceScriptQuery{ServiceID: &serviceID})
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
