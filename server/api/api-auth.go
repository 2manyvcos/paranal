package api

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/2manyvcos/paranal/server/helper"
)

func PostAuth(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var payload AuthRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil || payload.Username == "" || payload.Password == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := app.GetUser(payload.Username)
	if errors.Is(err, data.ErrNotFound) {
		maskAuthRejection()
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("Error loading user - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if user.PasswordHash == "" {
		maskAuthRejection()
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	matches, err := crypto.Argon2IDCompare(payload.Password, user.PasswordHash)
	if !matches {
		if err != nil {
			log.Printf("Error comparing hash - %s\n", err)
		}
		maskAuthRejection()
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	accessToken, expires, err := crypto.JWTGenerateToken(app.Config.Auth.JWT.Secret, "api", user.Name)
	if err != nil {
		log.Printf("Error generating access token - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(AuthResponse{
		Username:    payload.Username,
		AccessToken: accessToken,
		Expires:     expires,
	})
}

// Obscures the existence of the requested user and adds a random delay as a simple (but not too effective) brute force protection mechanism
func maskAuthRejection() {
	time.Sleep((2 * time.Second) + (time.Duration(rand.Intn(7000)) * time.Millisecond))
}
