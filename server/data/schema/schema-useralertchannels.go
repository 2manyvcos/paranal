package schema

import "fmt"

type UserAlertChannelQuery struct {
	ID       *int
	UserName *string
}

type UserAlertChannel struct {
	ID       int
	UserName string
	URL      string
}

func (r UserAlertChannel) Valid() error {
	if r.UserName == "" {
		return fmt.Errorf("invalid user name")
	}
	if r.URL == "" {
		return fmt.Errorf("invalid URL")
	}
	return nil
}
