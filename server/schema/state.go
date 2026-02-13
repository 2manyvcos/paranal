package schema

import (
	"time"
)

type State interface {
	ListServiceScriptStates(serviceID string) map[string]ServiceScriptState
	GetServiceScriptState(serviceID string, scriptID string) *ServiceScriptState
	RunServiceScripts()
	RunServiceScript(serviceID string, scriptID string) error
	OnServiceScriptChanged(script ServiceScript)
	OnServiceScriptDeleted(serviceId string, scriptID string)
	ListServiceActionGroups(serviceID string) []ServiceActionGroup
	ListServiceActions(serviceID string) []ServiceAction
	GetServiceAction(serviceID string, actionName string) *ServiceAction
	ListServiceUptimeStatuses(serviceID string) []ServiceUptimeStatus
	ListServiceVersions(serviceID string) []ServiceVersion
	GetServiceVersionDetails(serviceID string, versionName string) *ServiceVersionDetails
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

type ServiceUptimeStatus struct {
	Name      string
	Status    int
	Unhealthy bool
}

const (
	ServiceUptimeStatusDown = iota + 1
	ServiceUptimeStatusUp
)

var (
	ServiceUptimeStatusNames = map[int]string{
		ServiceUptimeStatusDown: "down",
		ServiceUptimeStatusUp:   "up",
	}
	ServiceUptimeStatusCodes map[string]int
)

func init() {
	ServiceUptimeStatusCodes = make(map[string]int, len(ServiceUptimeStatusNames))
	for code, name := range ServiceUptimeStatusNames {
		ServiceUptimeStatusCodes[name] = code
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
