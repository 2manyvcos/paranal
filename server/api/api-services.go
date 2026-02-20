package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
	"github.com/google/uuid"
)

func GetServices(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var query schema.UserServiceQuery
	q := req.URL.Query()
	if v, ok := utils.LoadQueryValue(q, "id"); ok {
		query.ID = &v
	}
	if v, ok := utils.LoadQueryValue(q, "name"); ok {
		query.Name = &v
	}
	if v, ok := utils.LoadQueryValue(q, "description"); ok {
		query.Description = &v
	}
	if v, ok := utils.LoadQueryValue(q, "logo"); ok {
		query.Logo = &v
	}
	if v, ok := utils.LoadQueryValue(q, "url"); ok {
		query.URL = &v
	}
	if v, ok := utils.LoadQueryBool(q, "config.favorite"); ok {
		query.Favorite = &v
	}
	if v, ok := utils.LoadQueryBool(q, "config.hidden"); ok {
		query.Hidden = &v
	}
	if v, ok := utils.LoadQueryBool(q, "config.uptimeAlerts"); ok {
		query.UptimeAlerts = &v
	}
	if v, ok := utils.LoadQueryBool(q, "config.versionAlerts"); ok {
		query.VersionAlerts = &v
	}

	records, err := app.ListUserServices(authorizedUser.Name, &query)
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
	app.PublishClientEvent(schema.NewUpdateEvent("/v1/services/" + url.PathEscape(newRecord.ID)))
	res.WriteHeader(http.StatusCreated)
}

func GetServicesByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUserService(
		authorizedUser.Name,
		schema.UserServiceQuery{
			ServiceQuery: schema.ServiceQuery{ID: &serviceID},
		},
	)
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
	app.PublishClientEvent(schema.NewUpdateEvent("/v1/services/" + url.PathEscape(serviceID)))
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
	app.OnServiceDeleted(serviceID)
	app.PublishClientEvent(schema.NewUpdateEvent("/v1/services/" + url.PathEscape(serviceID)))
}

func PatchServicesByIDConfig(res http.ResponseWriter, req *http.Request) {
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

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUserService(
		authorizedUser.Name,
		schema.UserServiceQuery{
			ServiceQuery: schema.ServiceQuery{ID: &serviceID},
		},
	)
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

	updatedRecord := schema.ServiceConfig{
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
	err = app.CreateOrUpdateServiceConfig(updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(authorizedUser.Name, "/v1/services/"+url.PathEscape(updatedRecord.ServiceID)))
}
