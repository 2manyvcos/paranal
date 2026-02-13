package schema

import (
	"time"
)

type State interface {
	GetServiceScriptStates(serviceID string) map[string]ServiceScriptState
	GetServiceScriptState(serviceID string, scriptID string) *ServiceScriptState
	GetServiceActionGroups(serviceID string) []ServiceActionGroup
	GetServiceActions(serviceID string) []ServiceAction
	GetServiceAction(serviceID string, actionName string) *ServiceAction
	RunServiceScripts()
	RunServiceScript(serviceID string, scriptID string) error
	OnServiceScriptChanged(script ServiceScript)
	OnServiceScriptDeleted(serviceId string, scriptID string)
	OnServiceDeleted(serviceID string)
}

type ServiceScriptState struct {
	LastRun *time.Time
	Running bool
	Error   error
	NextRun *time.Time
	RunNow  error
}

type ServiceActionGroup struct {
	Name string
	Icon string
}

type ServiceAction struct {
	Name             string
	Icon             string
	URL              string
	Script           string
	Group            string
	RestrictToAdmins bool
}
