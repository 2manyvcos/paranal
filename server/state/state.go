package state

import (
	"sync"
	"time"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/go-co-op/gocron/v2"
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
	lock                  sync.RWMutex
	scripts               map[string]*serviceScriptState
	uptimeStatuses        map[string]ServiceUptimeStatusState
	uptimeStatusesSorted  []ServiceUptimeStatusState
	versions              map[string]ServiceVersionState
	versionsSorted        []ServiceVersionState
	contextSections       map[string]ServiceContextSectionState
	contextSectionsSorted []ServiceContextSectionState
	contextOptions        map[string]ServiceContextOptionState
	contextOptionsSorted  []ServiceContextOptionState
}

type serviceScriptState struct {
	job             gocron.Job
	error           error
	uptimeStatuses  map[string]ServiceUptimeStatusState
	versions        map[string]ServiceVersionState
	contextSections map[string]ServiceContextSectionState
	contextOptions  map[string]ServiceContextOptionState
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

func (s ServiceUptimeStatusState) Up() bool {
	return s.Status == ServiceUptimeStatusUp
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
	Name                string `mapstructure:"name"`
	Time                time.Time
	Order               string              `mapstructure:"order"`
	CurrentVersion      string              `mapstructure:"currentVersion"`
	CurrentVersionNotes string              `mapstructure:"currentVersionNotes"`
	CurrentVersionCVEs  []ServiceVersionCVE `mapstructure:"currentVersionCVEs"`
	LatestVersion       string              `mapstructure:"latestVersion"`
	LatestVersionNotes  string              `mapstructure:"latestVersionNotes"`
	LatestVersionCVEs   []ServiceVersionCVE `mapstructure:"latestVersionCVEs"`
	Status              int
}

func (s ServiceVersionState) UpToDate() bool {
	return s.Status == ServiceVersionStatusUpToDate
}

func (s ServiceVersionState) Vulnerable() bool {
	return len(s.CurrentVersionCVEs) > 0
}

type ServiceVersionCVE struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	URL         string `mapstructure:"url"`
}

type ServiceContextSectionState struct {
	Name  string `mapstructure:"name"`
	Time  time.Time
	Order string `mapstructure:"order"`
	Icon  string `mapstructure:"icon"`
}

type ServiceContextOptionState struct {
	Name             string `mapstructure:"name"`
	Time             time.Time
	Order            string `mapstructure:"order"`
	Icon             string `mapstructure:"icon"`
	URL              string `mapstructure:"url"`
	Script           string `mapstructure:"script"`
	Section          string `mapstructure:"section"`
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

type ContextSectionInstruction struct {
	Type string `mapstructure:"type"`
	ServiceContextSectionState
	Time any `mapstructure:"time"`
}

type ContextOptionInstruction struct {
	Type string `mapstructure:"type"`
	ServiceContextOptionState
	Time any `mapstructure:"time"`
}
