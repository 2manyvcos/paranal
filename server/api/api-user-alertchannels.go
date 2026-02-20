package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
	"github.com/google/uuid"
)

func GetUserAlertChannels(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	query := schema.UserAlertChannelQuery{UserName: &authorizedUser.Name}
	q := req.URL.Query()
	if v, ok := utils.LoadQueryValue(q, "id"); ok {
		query.ID = &v
	}
	if v, ok := utils.LoadQueryValue(q, "userName"); ok {
		query.UserName = &v
	}
	if v, ok := utils.LoadQueryBool(q, "hasURL"); ok {
		query.HasURL = &v
	}
	if v, ok := utils.LoadQueryBool(q, "errorAlerts"); ok {
		query.ErrorAlerts = &v
	}
	if v, ok := utils.LoadQueryBool(q, "uptimeAlerts"); ok {
		query.UptimeAlerts = &v
	}
	if v, ok := utils.LoadQueryBool(q, "versionAlerts"); ok {
		query.VersionAlerts = &v
	}

	records, err := app.ListUserAlertChannels(&query)
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responsePayload := make([]struct {
		ID            string `json:"id"`
		URL           string `json:"url"`
		ErrorAlerts   bool   `json:"errorAlerts"`
		UptimeAlerts  bool   `json:"uptimeAlerts"`
		VersionAlerts bool   `json:"versionAlerts"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].ID = record.ID
		responsePayload[i].ErrorAlerts = record.ErrorAlerts
		responsePayload[i].UptimeAlerts = record.UptimeAlerts
		responsePayload[i].VersionAlerts = record.VersionAlerts
		if record.URL != "" {
			responsePayload[i].URL, err = crypto.Decrypt(app.Config.SecretKey, record.URL)
			if err != nil {
				log.Printf("Error decrypting value - %s\n", err)
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PostUserAlertChannels(res http.ResponseWriter, req *http.Request) {
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

	var requestPayload struct {
		URL           string `json:"url"`
		ErrorAlerts   bool   `json:"errorAlerts"`
		UptimeAlerts  bool   `json:"uptimeAlerts"`
		VersionAlerts bool   `json:"versionAlerts"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newRecord := schema.UserAlertChannel{
		ID:            uuid.NewString(),
		UserName:      authorizedUser.Name,
		ErrorAlerts:   requestPayload.ErrorAlerts,
		UptimeAlerts:  requestPayload.UptimeAlerts,
		VersionAlerts: requestPayload.VersionAlerts,
	}
	if requestPayload.URL != "" {
		newRecord.URL, err = crypto.Encrypt(app.Config.SecretKey, requestPayload.URL)
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
	err = app.CreateUserAlertChannel(newRecord)
	if errors.Is(err, schema.ErrConflict) {
		http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Error creating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(authorizedUser.Name, "/v1/user/alertchannels/"+url.PathEscape(newRecord.ID)))
	res.WriteHeader(http.StatusCreated)
}

func GetUserAlertChannelsByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	channelID := req.PathValue("channelID")
	if channelID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUserAlertChannel(schema.UserAlertChannelQuery{ID: &channelID, UserName: &authorizedUser.Name})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var url string
	if record.URL != "" {
		url, err = crypto.Decrypt(app.Config.SecretKey, record.URL)
		if err != nil {
			log.Printf("Error decrypting value - %s\n", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(struct {
		ID            string `json:"id"`
		URL           string `json:"url"`
		ErrorAlerts   bool   `json:"errorAlerts"`
		UptimeAlerts  bool   `json:"uptimeAlerts"`
		VersionAlerts bool   `json:"versionAlerts"`
	}{
		ID:            record.ID,
		URL:           url,
		ErrorAlerts:   record.ErrorAlerts,
		UptimeAlerts:  record.UptimeAlerts,
		VersionAlerts: record.VersionAlerts,
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PatchUserAlertChannelsByID(res http.ResponseWriter, req *http.Request) {
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

	channelID := req.PathValue("channelID")
	if channelID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUserAlertChannel(schema.UserAlertChannelQuery{ID: &channelID, UserName: &authorizedUser.Name})
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
		URL           utils.Optional[string] `json:"url"`
		ErrorAlerts   utils.Optional[bool]   `json:"errorAlerts"`
		UptimeAlerts  utils.Optional[bool]   `json:"uptimeAlerts"`
		VersionAlerts utils.Optional[bool]   `json:"versionAlerts"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	updatedRecord := record
	requestPayload.ErrorAlerts.ApplyIfDefined(&updatedRecord.ErrorAlerts)
	requestPayload.UptimeAlerts.ApplyIfDefined(&updatedRecord.UptimeAlerts)
	requestPayload.VersionAlerts.ApplyIfDefined(&updatedRecord.VersionAlerts)
	if requestPayload.URL.IsDefined {
		if requestPayload.URL.Value == "" {
			updatedRecord.URL = ""
		} else {
			updatedRecord.URL, err = crypto.Encrypt(app.Config.SecretKey, requestPayload.URL.Value)
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
	err = app.UpdateUserAlertChannels(schema.UserAlertChannelQuery{ID: &channelID, UserName: &authorizedUser.Name}, updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(authorizedUser.Name, "/v1/user/alertchannels/"+url.PathEscape(channelID)))
}

func DeleteUserAlertChannelsByID(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	channelID := req.PathValue("channelID")
	if channelID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := app.DeleteUserAlertChannels(schema.UserAlertChannelQuery{ID: &channelID, UserName: &authorizedUser.Name})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error deleting record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(authorizedUser.Name, "/v1/user/alertchannels/"+url.PathEscape(channelID)))
}
