package schema

import "fmt"

type ServiceVersionQuery struct {
	Name      *string
	ServiceID *string
}

type ServiceVersion struct {
	Name      string
	ServiceID string
	Version   string
}

func (r ServiceVersion) Valid() error {
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	if r.Version == "" {
		return fmt.Errorf("invalid version")
	}
	return nil
}
