package schema

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type ServiceScriptQuery struct {
	ID        *string
	ServiceID *string
}

type ServiceScript struct {
	ID        string
	Name      string
	Schedule  string
	Source    string
	ServiceID string
}

var serviceScriptCron = gocron.NewDefaultCron(true)

func (r ServiceScript) Valid() error {
	if r.ID == "" {
		return fmt.Errorf("invalid ID")
	}
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if r.Schedule == "" || serviceScriptCron.IsValid(r.Schedule, time.Local, time.Now()) != nil {
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
