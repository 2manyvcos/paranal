package application

import (
	"time"

	"github.com/2manyvcos/paranal/server/data/schema"
)

type State interface {
	GetServiceScriptStates(serviceID string) map[string]ServiceScriptState
	GetServiceScriptState(serviceID string, scriptID string) ServiceScriptState
	OnServiceScriptChanged(script schema.ServiceScript)
	OnServiceScriptDeleted(serviceId string, scriptID string)
	OnServiceDeleted(serviceID string)
}

type ServiceScriptState struct {
	LastRun *time.Time
	Error   error
	NextRun *time.Time
}
