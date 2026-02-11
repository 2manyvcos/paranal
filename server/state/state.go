package state

import (
	"log"
	"sync"
	"time"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/go-co-op/gocron/v2"
)

func Setup(app *application.App) error {
	s := State{app: app}

	scripts, err := app.ListServiceScripts(nil)
	if err != nil {
		return err
	}

	s.services = make(map[string]*serviceState)
	for _, script := range scripts {
		s.updateServiceScript(script)
	}

	app.State = &s
	app.Scheduler.Start()

	for serviceID, service := range s.services {
		for scriptID, script := range service.scripts {
			err := script.job.RunNow()
			if err != nil {
				log.Printf("Error running script \"%s\" for service \"%s\" - %s\n", scriptID, serviceID, err)
			}
		}
	}

	return nil
}

func (s *State) GetServiceScriptStates(serviceID string) map[string]application.ServiceScriptState {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return nil
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	result := make(map[string]application.ServiceScriptState, len(service.scripts))
	for scriptID, script := range service.scripts {
		var state application.ServiceScriptState
		if lastRun, err := script.job.LastRun(); err == nil && !lastRun.IsZero() {
			state.LastRun = &lastRun
		}
		state.Error = script.error
		if nextRun, err := script.job.NextRun(); err == nil && !nextRun.IsZero() {
			state.NextRun = &nextRun
		}
		result[scriptID] = state
	}
	return result
}

func (s *State) GetServiceScriptState(serviceID string, scriptID string) application.ServiceScriptState {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return application.ServiceScriptState{}
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	script, ok := service.scripts[scriptID]
	if !ok {
		return application.ServiceScriptState{}
	}
	var state application.ServiceScriptState
	if lastRun, err := script.job.LastRun(); err != nil && !lastRun.IsZero() {
		state.LastRun = &lastRun
	}
	state.Error = script.error
	if nextRun, err := script.job.NextRun(); err != nil && !nextRun.IsZero() {
		state.NextRun = &nextRun
	}
	return state
}

func (s *State) OnServiceScriptChanged(script schema.ServiceScript) {
	scriptState := s.updateServiceScript(script)
	if scriptState != nil {
		err := scriptState.job.RunNow()
		if err != nil {
			log.Printf("Error running script \"%s\" for service \"%s\" - %s\n", script.ID, script.ServiceID, err)
		}
	}
}

func (s *State) updateServiceScript(script schema.ServiceScript) *serviceScriptState {
	s.servicesLock.Lock()
	defer s.servicesLock.Unlock()
	service, ok := s.services[script.ServiceID]
	if !ok {
		service = &serviceState{
			scripts: make(map[string]*serviceScriptState),
		}
		s.services[script.ServiceID] = service
	}
	service.lock.Lock()
	defer service.lock.Unlock()
	if script, ok := service.scripts[script.ID]; ok {
		s.app.Scheduler.RemoveJob(script.job.ID())
	}
	var scriptState serviceScriptState
	var err error
	scriptState.job, err = s.app.Scheduler.NewJob(
		gocron.CronJob(script.Schedule, true),
		gocron.NewTask(newServiceScriptRunner(s, service, &scriptState, script)),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Printf("Error scheduling script \"%s\" for service \"%s\" - %s\n", script.ID, script.ServiceID, err)
		return nil
	}
	service.scripts[script.ID] = &scriptState
	return &scriptState
}

func (s *State) OnServiceScriptDeleted(serviceID string, scriptID string) {
	s.servicesLock.Lock()
	defer s.servicesLock.Unlock()
	service, ok := s.services[serviceID]
	if !ok {
		return
	}
	service.lock.Lock()
	defer service.lock.Unlock()
	if script, ok := service.scripts[scriptID]; ok {
		s.app.Scheduler.RemoveJob(script.job.ID())
		delete(service.scripts, scriptID)
	}
	if len(service.scripts) == 0 {
		delete(s.services, serviceID)
	} else {
		updateServiceState(service)
	}
}

func (s *State) OnServiceDeleted(serviceID string) {
	s.servicesLock.Lock()
	defer s.servicesLock.Unlock()
	if service, ok := s.services[serviceID]; ok {
		service.lock.Lock()
		for _, script := range service.scripts {
			s.app.Scheduler.RemoveJob(script.job.ID())
		}
		service.lock.Unlock()
		delete(s.services, serviceID)
	}
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
