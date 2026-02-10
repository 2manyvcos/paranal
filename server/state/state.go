package state

import (
	"log"
	"sync"
	"time"

	"github.com/2manyvcos/paranal/server/application"
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
		service, ok := s.services[script.ServiceID]
		if !ok {
			service = &serviceState{
				scripts: make(map[string]*serviceScriptState),
			}
			s.services[script.ServiceID] = service
		}
		if script, ok := service.scripts[script.ID]; ok {
			app.Scheduler.RemoveJob(script.job.ID())
		}
		var scriptState serviceScriptState
		scriptState.job, err = app.Scheduler.NewJob(
			gocron.CronJob(script.Schedule, true),
			gocron.NewTask(newServiceScriptRunner(&s, service, &scriptState, script)),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)
		if err != nil {
			log.Printf("Error scheduling script \"%s\" for service \"%s\" - %s\n", script.ID, script.ServiceID, err)
			continue
		}
		service.scripts[script.ID] = &scriptState
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
	lock                  sync.Mutex
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
