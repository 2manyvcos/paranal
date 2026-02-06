package schema

import "fmt"

type ServiceFavoriteQuery struct {
	UserName  *string
	ServiceID *string
}

type ServiceFavorite struct {
	UserName  string
	ServiceID string
}

func (r ServiceFavorite) Valid() error {
	if r.UserName == "" {
		return fmt.Errorf("invalid user name")
	}
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	return nil
}
