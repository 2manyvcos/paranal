package schema

import "fmt"

type UserAlertChannelQuery struct {
	ID            *string
	UserName      *string
	HasURL        *bool
	ErrorAlerts   *bool
	UptimeAlerts  *bool
	VersionAlerts *bool
}

type UserAlertChannel struct {
	ID            string
	UserName      string
	URL           string
	ErrorAlerts   bool
	UptimeAlerts  bool
	VersionAlerts bool
}

func (r UserAlertChannel) Valid() error {
	if r.ID == "" {
		return fmt.Errorf("invalid ID")
	}
	if r.UserName == "" {
		return fmt.Errorf("invalid user name")
	}
	if r.URL == "" {
		return fmt.Errorf("invalid URL")
	}
	return nil
}
