package schema

import (
	"context"
	"time"
)

type State interface {
	Close()

	PublishClientEvent(event ClientEvent)
	SubscribeToClientEvents(context context.Context) <-chan ClientEvent

	ListServiceScriptStates(serviceID string) map[string]ServiceScriptState
	GetServiceScriptState(serviceID string, scriptID string) *ServiceScriptState
	RunServiceScripts()
	RunServiceScript(serviceID string, scriptID string) error
	OnServiceScriptChanged(script ServiceScript)
	OnServiceScriptDeleted(serviceId string, scriptID string)
	ListServiceActionGroups(serviceID string) []ServiceActionGroup
	ListServiceActions(serviceID string) []ServiceAction
	GetServiceAction(serviceID string, actionName string) *ServiceAction
	ListServiceHealthStatuses(serviceID string) []ServiceHealthStatus
	ListServiceVersions(serviceID string) []ServiceVersion
	GetServiceVersionDetails(serviceID string, versionName string) *ServiceVersionDetails
	OnServiceDeleted(serviceID string)

	ListMaintenanceTaskStates() []MaintenanceTaskState
	GetMaintenanceTaskState(taskName string) *MaintenanceTaskState
	RunMaintenanceTask(taskName string) error
}

type ClientEvent interface {
	Event() string
	User() string
	Data() string
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

type ServiceHealthStatus struct {
	Name      string
	Status    int
	Unhealthy bool
}

const (
	ServiceHealthStatusDown = iota + 1
	ServiceHealthStatusUp
)

var (
	ServiceHealthStatusNames = map[int]string{
		ServiceHealthStatusDown: "down",
		ServiceHealthStatusUp:   "up",
	}
	ServiceHealthStatusCodes map[string]int
)

func init() {
	ServiceHealthStatusCodes = make(map[string]int, len(ServiceHealthStatusNames))
	for code, name := range ServiceHealthStatusNames {
		ServiceHealthStatusCodes[name] = code
	}
}

type ServiceVersion struct {
	Name           string
	CurrentVersion string
	CurrentCVEs    int
	LatestVersion  string
	LatestCVEs     int
	Status         int
	Outdated       bool
	Vulnerable     bool
}

const (
	ServiceVersionStatusOutdated = iota + 1
	ServiceVersionStatusUpToDate
)

var (
	ServiceVersionStatusNames = map[int]string{
		ServiceVersionStatusOutdated: "outdated",
		ServiceVersionStatusUpToDate: "upToDate",
	}
	ServiceVersionStatusCodes map[string]int
)

func init() {
	ServiceVersionStatusCodes = make(map[string]int, len(ServiceVersionStatusNames))
	for code, name := range ServiceVersionStatusNames {
		ServiceVersionStatusCodes[name] = code
	}
}

type ServiceVersionDetails struct {
	CurrentVersionNotes    string
	CurrentCVEDescriptions []ServiceCVE
	LatestVersionNotes     string
	LatestCVEDescriptions  []ServiceCVE
}

type ServiceCVE struct {
	Name        string
	Description string
	URL         string
}

type MaintenanceTaskState struct {
	Name    string
	LastRun *time.Time
	Running bool
	NextRun *time.Time
}
