package schema

import "fmt"

type FavoriteQuery struct {
	UserName  *string
	ServiceID *string
}

type Favorite struct {
	UserName  string
	ServiceID string
}

func (r Favorite) Valid() error {
	if r.UserName == "" {
		return fmt.Errorf("invalid user name")
	}
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	return nil
}
