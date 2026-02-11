package state

import (
	"sync"
	"time"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/go-co-op/gocron/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
)

func Setup(app *application.App) error {
	s := State{app: app}

	s.setupServiceScripts()

	return nil
}

type State struct {
	app          *application.App
	services     map[string]*serviceState
	servicesLock sync.RWMutex
}

type serviceState struct {
	lock                 sync.RWMutex
	scripts              map[string]*serviceScriptState
	uptimeStatuses       map[string]ServiceUptimeStatusState
	uptimeStatusesSorted []ServiceUptimeStatusState
	versions             map[string]ServiceVersionState
	versionsSorted       []ServiceVersionState
	actionGroups         map[string]ServiceActionGroupState
	actionGroupsSorted   []ServiceActionGroupState
	actions              map[string]ServiceActionState
	actionsSorted        []ServiceActionState
}

type serviceScriptState struct {
	program        jpl.JPLProgram
	job            gocron.Job
	error          error
	uptimeStatuses map[string]ServiceUptimeStatusState
	versions       map[string]ServiceVersionState
	actionGroups   map[string]ServiceActionGroupState
	actions        map[string]ServiceActionState
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

type ServiceUptimeStatusState struct {
	Name   string `mapstructure:"name"`
	Time   time.Time
	Order  string `mapstructure:"order"`
	Status int
}

func (s ServiceUptimeStatusState) Unhealthy() bool {
	return s.Status != ServiceUptimeStatusUp
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

type ServiceVersionState struct {
	Name                   string `mapstructure:"name"`
	Time                   time.Time
	Order                  string       `mapstructure:"order"`
	CurrentVersion         string       `mapstructure:"currentVersion"`
	CurrentVersionNotes    string       `mapstructure:"currentVersionNotes"`
	CurrentCVEs            int          `mapstructure:"currentCVEs"`
	CurrentCVEDescriptions []ServiceCVE `mapstructure:"currentCVEDescriptions"`
	LatestVersion          string       `mapstructure:"latestVersion"`
	LatestVersionNotes     string       `mapstructure:"latestVersionNotes"`
	LatestCVEs             int          `mapstructure:"latestCVEs"`
	LatestCVEDescriptions  []ServiceCVE `mapstructure:"latestCVEDescriptions"`
	Status                 int
}

func (s ServiceVersionState) Outdated() bool {
	return s.Status != ServiceVersionStatusUpToDate
}

func (s ServiceVersionState) Vulnerable() bool {
	return s.CurrentCVEs > 0
}

type ServiceCVE struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	URL         string `mapstructure:"url"`
}

type ServiceActionGroupState struct {
	Name  string `mapstructure:"name"`
	Time  time.Time
	Order string `mapstructure:"order"`
	Icon  string `mapstructure:"icon"`
}

type ServiceActionState struct {
	Name             string `mapstructure:"name"`
	Time             time.Time
	Order            string `mapstructure:"order"`
	Icon             string `mapstructure:"icon"`
	URL              string `mapstructure:"url"`
	Script           string `mapstructure:"script"`
	Group            string `mapstructure:"group"`
	RestrictToAdmins bool   `mapstructure:"restrictToAdmins"`
}

type GenericInstruction struct {
	Type string `mapstructure:"type"`
}

type UptimeStatusInstruction struct {
	Type string `mapstructure:"type"`
	ServiceUptimeStatusState
	Time   any    `mapstructure:"time"`
	Status string `mapstructure:"status"`
}

type VersionInstruction struct {
	Type string `mapstructure:"type"`
	ServiceVersionState
	Time   any    `mapstructure:"time"`
	Status string `mapstructure:"status"`
}

type ActionGroupInstruction struct {
	Type string `mapstructure:"type"`
	ServiceActionGroupState
	Time any `mapstructure:"time"`
}

type ActionInstruction struct {
	Type string `mapstructure:"type"`
	ServiceActionState
	Time any `mapstructure:"time"`
}
