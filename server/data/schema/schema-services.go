package schema

import "fmt"

type ServiceQuery struct {
	ID *string
}

type Service struct {
	ID          string
	Name        string
	Description string
	Logo        string
	URL         string
}

func (r Service) Valid() error {
	if r.ID == "" {
		return fmt.Errorf("invalid ID")
	}
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

type ServiceWithFavorite struct {
	Service
	IsFavorite bool
}
