package schema

import "fmt"

type UserAlertChannelQuery struct {
	ID       *string
	UserName *string
}

type UserAlertChannel struct {
	ID       string
	UserName string
	URL      string
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
