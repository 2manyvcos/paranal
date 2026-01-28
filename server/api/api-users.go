package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

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

	result := make([]User, len(users))
	for i, user := range users {
		result[i] = UserFromData(user)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(result)
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
	json.NewEncoder(res).Encode(UserFromData(user))
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
