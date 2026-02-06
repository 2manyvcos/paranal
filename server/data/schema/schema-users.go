package schema

import (
	"fmt"
)

type UserQuery struct {
	Name          *string
	ErrorAlerts   *bool
	UptimeAlerts  *bool
	VersionAlerts *bool
}

type User struct {
	Name          string
	DisplayName   string
	Role          int
	PasswordHash  string
	ErrorAlerts   bool
	UptimeAlerts  bool
	VersionAlerts bool
}

const (
	UserRoleCommon = iota + 1
	UserRoleAdmin
)

var (
	UserRoleNames = map[int]string{
		UserRoleCommon: "common",
		UserRoleAdmin:  "admin",
	}
	UserRoleCodes map[string]int
)

func init() {
	UserRoleCodes = make(map[string]int, len(UserRoleNames))
	for code, name := range UserRoleNames {
		UserRoleCodes[name] = code
	}
}

func (r User) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if _, ok := UserRoleNames[r.Role]; !ok {
		return fmt.Errorf("invalid role")
	}
	return nil
}
