package schema

import "fmt"

type ServiceUserConfigQuery struct {
	UserName     *string
	ServiceID    *string
	Favorite     *bool
	UptimeAlert  *bool
	VersionAlert *bool
}

type ServiceUserConfig struct {
	UserName     string
	ServiceID    string
	Favorite     bool
	UptimeAlert  bool
	VersionAlert bool
}

func (r ServiceUserConfig) Valid() error {
	if r.UserName == "" {
		return fmt.Errorf("invalid user name")
	}
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	return nil
}
