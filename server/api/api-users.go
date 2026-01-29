package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/2manyvcos/paranal/server/helper"
)

func GetUsers(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	users, err := app.ListUsers()
	if err != nil {
		log.Printf("Error loading users - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	result := make([]struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Role        string `json:"role"`
	}, len(users))
	for i, user := range users {
		result[i].Name = user.Name
		result[i].DisplayName = user.DisplayName
		result[i].Role = data.USER_ROLE_NAMES[user.Role]
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(result)
}

func PostUsers(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var payload struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Role        string `json:"role"`
		Password    string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newUser := data.User{
		Name:        payload.Name,
		DisplayName: payload.DisplayName,
		Role:        data.USER_ROLE_CODES[payload.Role],
	}
	if payload.Password != "" {
		newUser.PasswordHash, err = crypto.Hash(payload.Password)
		if err != nil {
			log.Printf("Error while generating hash - %s\n", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if !newUser.Valid() {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.CreateUser(newUser, false)
	if errors.Is(err, data.ErrConflict) {
		http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Error creating user - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
}

func GetUsersByUsername(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	username := req.PathValue("username")

	user, err := app.GetUser(username)
	if errors.Is(err, data.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading user - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Role        string `json:"role"`
	}{
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Role:        data.USER_ROLE_NAMES[user.Role],
	})
}

func PutUsersByUsername(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	username := req.PathValue("username")

	if authorizedUser != nil && authorizedUser.Name == username {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var payload struct {
		DisplayName string `json:"displayName"`
		Role        string `json:"role"`
		Password    string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := app.GetUser(username)
	if errors.Is(err, data.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error loading user - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newUser := user
	newUser.DisplayName = payload.DisplayName
	newUser.Role = data.USER_ROLE_CODES[payload.Role]
	if payload.Password != "" {
		newUser.PasswordHash, err = crypto.Hash(payload.Password)
		if err != nil {
			log.Printf("Error while generating hash - %s\n", err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if !newUser.Valid() {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = app.UpdateUser(newUser)
	if err != nil {
		log.Printf("Error creating user - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func DeleteUsersByUsername(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	username := req.PathValue("username")

	if authorizedUser != nil && authorizedUser.Name == username {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err := app.DeleteUser(username)
	if errors.Is(err, data.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error deleting user - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
