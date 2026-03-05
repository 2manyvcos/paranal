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
)

func GetUsers(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	var query schema.UserQuery
	q := req.URL.Query()
	if v, ok := utils.LoadQueryValue(q, "name"); ok {
		query.Name = &v
	}
	if v, ok := utils.LoadQueryValue(q, "role"); ok {
		t, ok := schema.UserRoleCodes[v]
		if !ok {
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		query.Role = &t
	}
	if v, ok := utils.LoadQueryBool(q, "hasPassword"); ok {
		query.HasPassword = &v
	}

	records, err := app.ListUsers(&query)
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responsePayload := make([]struct {
		Name        string `json:"name"`
		Role        string `json:"role"`
		HasPassword bool   `json:"hasPassword"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].Name = record.Name
		responsePayload[i].Role = schema.UserRoleNames[record.Role]
		responsePayload[i].HasPassword = record.PasswordHash != ""
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PostUsers(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Name     string `json:"name"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newRecord := schema.User{
		Name: requestPayload.Name,
		Role: schema.UserRoleCodes[requestPayload.Role],
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
	err = app.CreateUser(newRecord)
	if errors.Is(err, schema.ErrConflict) {
		http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Error creating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(newRecord.Name, "/v1/user"))
	app.PublishClientEvent(schema.NewUpdateEvent("/v1/users/" + url.PathEscape(newRecord.Name)))

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(res).Encode(struct {
		Name string `json:"name"`
	}{
		Name: newRecord.Name,
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func GetUsersByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	userName := req.PathValue("userName")
	if userName == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUser(schema.UserQuery{Name: &userName})
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
		Name        string `json:"name"`
		Role        string `json:"role"`
		HasPassword bool   `json:"hasPassword"`
	}{
		Name:        record.Name,
		Role:        schema.UserRoleNames[record.Role],
		HasPassword: record.PasswordHash != "",
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PatchUsersByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	userName := req.PathValue("userName")
	if userName == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if authorizedUser != nil && authorizedUser.Name == userName {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Role     utils.Optional[string] `json:"role"`
		Password utils.Optional[string] `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUser(schema.UserQuery{Name: &userName})
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
	if requestPayload.Role.IsDefined {
		updatedRecord.Role = schema.UserRoleCodes[requestPayload.Role.Value]
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
	err = app.UpdateUsers(schema.UserQuery{Name: &userName}, updatedRecord)
	if err != nil {
		log.Printf("Error updating record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(userName, "/v1/user"))
	app.PublishClientEvent(schema.NewUpdateEvent("/v1/users/" + url.PathEscape(userName)))
}

func DeleteUsersByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	userName := req.PathValue("userName")
	if userName == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if authorizedUser != nil && authorizedUser.Name == userName {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err := app.DeleteUsers(schema.UserQuery{Name: &userName})
	if errors.Is(err, schema.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error deleting record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.PublishClientEvent(schema.NewUserUpdateEvent(userName, "/v1/user"))
	app.PublishClientEvent(schema.NewUpdateEvent("/v1/users/" + url.PathEscape(userName)))
}
