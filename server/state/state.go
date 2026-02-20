package state

import (
	"sync"
	"time"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
	"github.com/go-co-op/gocron/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
)

func Setup(app *application.App) (*State, error) {
	s := State{app: app}

	if err := s.setupClientEvents(); err != nil {
		return nil, err
	}

	if err := s.setupServiceScripts(); err != nil {
		return nil, err
	}

	if err := s.setupMaintenanceTasks(); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *State) Close() {
	s.closeClientEvents()
}

type State struct {
	app              *application.App
	clientEvents     *utils.Broker[schema.ClientEvent]
	services         map[string]*serviceState
	servicesLock     sync.RWMutex
	maintenanceTasks map[string]maintenanceTask
}

type serviceState struct {
	lock                 sync.RWMutex
	scripts              map[string]*serviceScriptState
	uptimeStatuses       map[string]ServiceUptimeStatus
	uptimeStatusesSorted []ServiceUptimeStatus
	versions             map[string]ServiceVersion
	versionsSorted       []ServiceVersion
	actionGroups         map[string]ServiceActionGroup
	actionGroupsSorted   []ServiceActionGroup
	actions              map[string]ServiceAction
	actionsSorted        []ServiceAction
}

type serviceScriptState struct {
	program        jpl.JPLProgram
	job            gocron.Job
	running        bool
	error          error
	uptimeStatuses map[string]ServiceUptimeStatus
	versions       map[string]ServiceVersion
	actionGroups   map[string]ServiceActionGroup
	actions        map[string]ServiceAction
}

type ServiceUptimeStatus struct {
	Name   string `mapstructure:"name"`
	Time   time.Time
	Order  string `mapstructure:"order"`
	Status int
}

func (s ServiceUptimeStatus) Unhealthy() bool {
	return s.Status != schema.ServiceUptimeStatusUp
}

type ServiceVersion struct {
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

func (s ServiceVersion) Outdated() bool {
	return s.Status != schema.ServiceVersionStatusUpToDate
}

func (s ServiceVersion) Vulnerable() bool {
	return s.CurrentCVEs > 0
}

type ServiceCVE struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	URL         string `mapstructure:"url"`
}

type ServiceActionGroup struct {
	Name  string `mapstructure:"name"`
	Time  time.Time
	Order string `mapstructure:"order"`
	Icon  string `mapstructure:"icon"`
}

type ServiceAction struct {
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
	ServiceUptimeStatus
	Time   any    `mapstructure:"time"`
	Status string `mapstructure:"status"`
}

type VersionInstruction struct {
	Type string `mapstructure:"type"`
	ServiceVersion
	Time   any    `mapstructure:"time"`
	Status string `mapstructure:"status"`
}

type ActionGroupInstruction struct {
	Type string `mapstructure:"type"`
	ServiceActionGroup
	Time any `mapstructure:"time"`
}

type ActionInstruction struct {
	Type string `mapstructure:"type"`
	ServiceAction
	Time any `mapstructure:"time"`
}

type maintenanceTask interface {
	State() (lastRun *time.Time, running bool, nextRun *time.Time)
	Run() error
}
