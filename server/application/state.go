package application

import (
	"time"

	"github.com/2manyvcos/paranal/server/data/schema"
)

type State interface {
	GetServiceScriptStates(serviceID string) map[string]ServiceScriptState
	GetServiceScriptState(serviceID string, scriptID string) *ServiceScriptState
	GetServiceActionGroups(serviceID string) []ServiceActionGroup
	GetServiceActions(serviceID string) []ServiceAction
	GetServiceAction(serviceID string, actionName string) *ServiceAction
	RunAllServiceScripts()
	RunServiceScript(serviceID string, scriptID string) (alreadyRunning bool, found bool)
	OnServiceScriptChanged(script schema.ServiceScript)
	OnServiceScriptDeleted(serviceId string, scriptID string)
	OnServiceDeleted(serviceID string)
}

type ServiceScriptState struct {
	LastRun *time.Time
	Running bool
	Error   error
	NextRun *time.Time
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
