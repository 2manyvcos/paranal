package schema

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type ServiceScriptQuery struct {
	ID        *string
	ServiceID *string
	Name      *string
	Schedule  *string
}

type ServiceScript struct {
	ID        string
	ServiceID string
	Name      string
	Schedule  string
	Source    string
}

var serviceScriptCron = gocron.NewDefaultCron(true)

func (r ServiceScript) Valid() error {
	if r.ID == "" {
		return fmt.Errorf("invalid ID")
	}
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
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
	return nil
}
