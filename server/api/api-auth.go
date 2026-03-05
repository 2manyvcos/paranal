package api

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
)

func PostAuth(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil || requestPayload.UserName == "" || requestPayload.Password == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.GetUser(schema.UserQuery{Name: &requestPayload.UserName})
	if errors.Is(err, schema.ErrNotFound) {
		maskAuthRejection()
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if record.PasswordHash == "" {
		maskAuthRejection()
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	matches, err := crypto.CompareToHash(requestPayload.Password, record.PasswordHash)
	if !matches {
		if err != nil {
			log.Printf("Error comparing hash - %s\n", err)
		}
		maskAuthRejection()
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	accessToken, expires, err := crypto.GenerateJWTToken(app.Config.Auth.JWT.Secret, "api", record.Name)
	if err != nil {
		log.Printf("Error generating access token - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(struct {
		UserName    string    `json:"userName"`
		AccessToken string    `json:"accessToken"`
		Expires     time.Time `json:"expires"`
	}{
		UserName:    requestPayload.UserName,
		AccessToken: accessToken,
		Expires:     expires,
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

// Obscures the existence of the requested user and adds a random delay as a simple (but not too effective) brute force protection mechanism
func maskAuthRejection() {
	time.Sleep((2 * time.Second) + (time.Duration(rand.Intn(7000)) * time.Millisecond))
}
