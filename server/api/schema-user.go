package api

import "github.com/2manyvcos/paranal/server/data"

type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
}

func UserFromData(other data.User) (result User) {
	result.Name = other.Name
	result.DisplayName = other.DisplayName
	result.Role = map[int]string{
		data.USER_ROLE_ADMIN:  "admin",
		data.USER_ROLE_COMMON: "common",
	}[other.Role]
	return
}
