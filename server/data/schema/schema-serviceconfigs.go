package schema

import "fmt"

type ServiceConfigQuery struct {
	UserName      *string
	ServiceID     *string
	Favorite      *bool
	Hidden        *bool
	UptimeAlerts  *bool
	VersionAlerts *bool
}

type ServiceConfig struct {
	UserName      string
	ServiceID     string
	Favorite      bool
	Hidden        bool
	UptimeAlerts  bool
	VersionAlerts bool
}

func (r ServiceConfig) Valid() error {
	if r.UserName == "" {
		return fmt.Errorf("invalid user name")
	}
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	return nil
}
