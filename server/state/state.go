package state

import (
	"log"
	"sync"
	"time"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/2manyvcos/paranal/server/scripts"
	"github.com/go-co-op/gocron/v2"
)

func Setup(app *application.App) error {
	s := State{app: app}

	lastKnownUptimeStatuses, err := app.ListServiceUptimeStatuses(nil)
	if err != nil {
		return err
	}
	lastKnownUptimeStatusMap := make(map[string]map[string]ServiceUptimeStatusState)
	for _, uptimeStatus := range lastKnownUptimeStatuses {
		service, ok := lastKnownUptimeStatusMap[uptimeStatus.ServiceID]
		if !ok {
			service = make(map[string]ServiceUptimeStatusState)
			lastKnownUptimeStatusMap[uptimeStatus.ServiceID] = service
		}
		service[uptimeStatus.Name] = ServiceUptimeStatusState{
			Name:   uptimeStatus.Name,
			Order:  uptimeStatus.Name,
			Status: uptimeStatus.Status,
			hidden: true,
		}
	}

	lastKnownVersions, err := app.ListServiceVersions(nil)
	if err != nil {
		return err
	}
	lastKnownVersionMap := make(map[string]map[string]ServiceVersionState)
	for _, version := range lastKnownVersions {
		service, ok := lastKnownVersionMap[version.ServiceID]
		if !ok {
			service = make(map[string]ServiceVersionState)
			lastKnownVersionMap[version.ServiceID] = service
		}
		service[version.Name] = ServiceVersionState{
			Name:          version.Name,
			Order:         version.Name,
			LatestVersion: version.Version,
			hidden:        true,
		}
	}

	scripts, err := app.ListServiceScripts(nil)
	if err != nil {
		return err
	}

	s.services = make(map[string]*serviceState)
	for _, script := range scripts {
		service, ok := s.services[script.ServiceID]
		if !ok {
			service = &serviceState{
				scripts: make(map[string]serviceScriptState),
			}
			s.services[script.ServiceID] = service
		}
		if script, ok := service.scripts[script.ID]; ok {
			app.Scheduler.RemoveJob(script.job.ID())
		}
		scriptState := serviceScriptState{
			uptimeStatuses: lastKnownUptimeStatusMap[script.ServiceID],
			versions:       lastKnownVersionMap[script.ServiceID],
		}
		scriptState.job, err = app.Scheduler.NewJob(
			gocron.CronJob(script.Schedule, true),
			gocron.NewTask(runScript, &s, service, &scriptState, script),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)
		if err != nil {
			log.Printf("Error scheduling script \"%s\" for service \"%s\" - %s\n", script.ID, script.ServiceID, err)
			continue
		}
		service.scripts[script.ID] = scriptState
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

type State struct {
	app          *application.App
	services     map[string]*serviceState
	servicesLock sync.Mutex
}

type serviceState struct {
	lock            sync.Mutex
	scripts         map[string]serviceScriptState
	uptimeStatuses  []ServiceUptimeStatusState
	versions        []ServiceVersionState
	contextSections []ServiceContextSectionState
	contextOptions  []ServiceContextOptionState
}

type serviceScriptState struct {
	job             gocron.Job
	uptimeStatuses  map[string]ServiceUptimeStatusState
	versions        map[string]ServiceVersionState
	contextSections map[string]ServiceContextSectionState
	contextOptions  map[string]ServiceContextOptionState
}

type ServiceUptimeStatusState struct {
	Name   string    `json:"name"`
	Time   time.Time `json:"time"`
	Order  string    `json:"order"`
	Status int
	hidden bool
}

type ServiceVersionState struct {
	Name                string              `json:"name"`
	Time                time.Time           `json:"time"`
	Order               string              `json:"order"`
	CurrentVersion      string              `json:"currentVersion"`
	CurrentVersionNotes string              `json:"currentVersionNotes"`
	CurrentVersionCVEs  []ServiceVersionCVE `json:"currentVersionCVEs"`
	LatestVersion       string              `json:"latestVersion"`
	LatestVersionNotes  string              `json:"latestVersionNotes"`
	LatestVersionCVEs   []ServiceVersionCVE `json:"latestVersionCVEs"`
	hidden              bool
}

type ServiceVersionCVE struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type ServiceContextSectionState struct {
	Name  string    `json:"name"`
	Time  time.Time `json:"time"`
	Order string    `json:"order"`
	Icon  string    `json:"icon"`
}

type ServiceContextOptionState struct {
	Name             string    `json:"name"`
	Time             time.Time `json:"time"`
	Order            string    `json:"order"`
	Icon             string    `json:"icon"`
	URL              string    `json:"url"`
	Script           string    `json:"script"`
	RestrictToAdmins bool      `json:"restrictToAdmins"`
}

func runScript(s *State, serviceState *serviceState, scriptState *serviceScriptState, script schema.ServiceScript) {
	results, err := scripts.RunServiceScript(s.app, script.ServiceID, script.Source)
	if err != nil {
		log.Printf("Script error [%s:%s]: %s", script.ServiceID, script.ID, err)
		return
	}
	log.Printf("Script result [%s:%s]: %+v", script.ServiceID, script.ID, results)
	// TODO:
	/*
	  - handle script results (and errors)
	  - integrate state into the rest api (fetching; update on data change)
	*/
}
