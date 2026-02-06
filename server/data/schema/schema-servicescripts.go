package schema

import "fmt"

type ServiceScriptQuery struct {
	ID        *int
	ServiceID *string
}

type ServiceScript struct {
	ID        int
	Name      string
	Schedule  string
	Source    string
	ServiceID string
}

func (r ServiceScript) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if r.Schedule == "" {
		return fmt.Errorf("invalid schedule")
	}
	if r.Source == "" {
		return fmt.Errorf("invalid source")
	}
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	return nil
}
