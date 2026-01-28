package api

import "time"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Username    string    `json:"username"`
	AccessToken string    `json:"accessToken"`
	Expires     time.Time `json:"expires"`
}
