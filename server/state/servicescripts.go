package state

import (
	"fmt"
	"log"
	"sort"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/2manyvcos/paranal/server/scripts"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-viper/mapstructure/v2"
)

func (s *State) setupServiceScripts() error {
	scripts, err := s.app.ListServiceScripts(nil)
	if err != nil {
		return err
	}

	s.services = make(map[string]*serviceState)
	for _, script := range scripts {
		s.updateServiceScript(script)
	}

	s.app.State = s
	s.app.Scheduler.Start()

	for serviceID, service := range s.services {
		for scriptID, script := range service.scripts {
			if script.job != nil {
				err := script.job.RunNow()
				if err != nil {
					log.Printf("Error running script \"%s\" for service \"%s\" - %s\n", scriptID, serviceID, err)
				}
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
		state := application.ServiceScriptState{
			Error: script.error,
		}
		if script.job != nil {
			if lastRun, err := script.job.LastRun(); err == nil && !lastRun.IsZero() {
				state.LastRun = &lastRun
			}
			if nextRun, err := script.job.NextRun(); err == nil && !nextRun.IsZero() {
				state.NextRun = &nextRun
			}
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
	state := application.ServiceScriptState{
		Error: script.error,
	}
	if script.job != nil {
		if lastRun, err := script.job.LastRun(); err != nil && !lastRun.IsZero() {
			state.LastRun = &lastRun
		}
		if nextRun, err := script.job.NextRun(); err != nil && !nextRun.IsZero() {
			state.NextRun = &nextRun
		}
	}
	return state
}

func (s *State) OnServiceScriptChanged(script schema.ServiceScript) {
	scriptState := s.updateServiceScript(script)
	if scriptState.job != nil {
		err := scriptState.job.RunNow()
		if err != nil {
			log.Printf("Error running script \"%s\" for service \"%s\" - %s\n", script.ID, script.ServiceID, err)
		}
	}
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
		if script.job != nil {
			s.app.Scheduler.RemoveJob(script.job.ID())
		}
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
			if script.job != nil {
				s.app.Scheduler.RemoveJob(script.job.ID())
			}
		}
		service.lock.Unlock()
		delete(s.services, serviceID)
	}
}

func (s *State) RunServiceScript(serviceID string, scriptID string) bool {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return false
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	script, ok := service.scripts[scriptID]
	if !ok {
		return false
	}
	if script.job != nil {
		err := script.job.RunNow()
		if err != nil {
			log.Printf("Error running script \"%s\" for service \"%s\" - %s\n", scriptID, serviceID, err)
		}
	}
	return true
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
	if script, ok := service.scripts[script.ID]; ok && script.job != nil {
		s.app.Scheduler.RemoveJob(script.job.ID())
	}
	var scriptState serviceScriptState
	if program, err := scripts.ParseServiceScript(s.app, script.ServiceID, script.Source); err == nil {
		scriptState.program = program
	} else {
		log.Printf("Error parsing script \"%s\" for service \"%s\" - %s\n", script.ID, script.ServiceID, err)
		scriptState.error = err
		alertScriptError(s.app, nil, script, err)
	}
	if scriptState.program != nil {
		if job, err := s.app.Scheduler.NewJob(
			gocron.CronJob(script.Schedule, true),
			gocron.NewTask(newServiceScriptRunner(s, service, &scriptState, script)),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		); err == nil {
			scriptState.job = job
		} else {
			log.Printf("Error scheduling script \"%s\" for service \"%s\" - %s\n", script.ID, script.ServiceID, err)
			scriptState.error = err
			alertScriptError(s.app, nil, script, err)
		}
	}
	service.scripts[script.ID] = &scriptState
	return &scriptState
}

func newServiceScriptRunner(s *State, serviceState *serviceState, scriptState *serviceScriptState, script schema.ServiceScript) func() {
	handleErr := func(service *schema.Service, err error) {
		log.Printf("Error running script \"%s\" for service \"%s\" - %s", script.ServiceID, script.ID, err)
		serviceState.lock.Lock()
		scriptState.error = err
		scriptState.uptimeStatuses = nil
		scriptState.versions = nil
		scriptState.actionGroups = nil
		scriptState.actions = nil
		updateServiceState(serviceState)
		serviceState.lock.Unlock()
		alertScriptError(s.app, service, script, err)
	}

	return func() {
		service, err := s.app.GetService(schema.ServiceQuery{ID: &script.ServiceID})
		if err != nil {
			handleErr(nil, err)
			return
		}

		results, err := scriptState.program.Run([]any{nil}, nil)
		if err != nil {
			handleErr(&service, err)
			return
		}

		var uptimeStatuses []ServiceUptimeStatusState
		var versions []ServiceVersionState
		var actionGroups []ServiceActionGroupState
		var actions []ServiceActionState

		for _, result := range results {
			var i GenericInstruction
			err := mapstructure.Decode(result, &i)
			if err != nil {
				handleErr(&service, err)
				return
			}

			switch i.Type {
			case "uptimeStatus":
				var i UptimeStatusInstruction
				if err = decodeServiceInstruction(result, &i); err != nil {
					handleErr(&service, fmt.Errorf("invalid instruction - %s", err))
					return
				}
				uptimeStatus := i.ServiceUptimeStatusState
				if uptimeStatus.Name == "" {
					uptimeStatus.Name = service.Name
				}
				if uptimeStatus.Time, err = decodeTime(i.Time); err != nil {
					handleErr(&service, err)
					return
				}
				if uptimeStatus.Order == "" {
					uptimeStatus.Order = uptimeStatus.Name
				}
				if status, ok := ServiceUptimeStatusCodes[i.Status]; ok {
					uptimeStatus.Status = status
				} else {
					handleErr(&service, fmt.Errorf("invalid status"))
					return
				}
				uptimeStatuses = append(uptimeStatuses, uptimeStatus)

			case "version":
				var i VersionInstruction
				if err = decodeServiceInstruction(result, &i); err != nil {
					handleErr(&service, fmt.Errorf("invalid instruction - %s", err))
					return
				}
				version := i.ServiceVersionState
				if version.Name == "" {
					version.Name = service.Name
				}
				if version.Time, err = decodeTime(i.Time); err != nil {
					handleErr(&service, err)
					return
				}
				if version.Order == "" {
					version.Order = version.Name
				}
				if version.CurrentVersion == "" {
					handleErr(&service, fmt.Errorf("invalid version"))
					return
				}
				if i.Status == "" {
					if version.LatestVersion == "" {
						handleErr(&service, fmt.Errorf("invalid version"))
						return
					}
					if version.CurrentVersion == version.LatestVersion {
						version.Status = ServiceVersionStatusUpToDate
					} else {
						version.Status = ServiceVersionStatusOutdated
					}
				} else if status, ok := ServiceVersionStatusCodes[i.Status]; ok {
					version.Status = status
				} else {
					handleErr(&service, fmt.Errorf("invalid status"))
					return
				}
				for _, cve := range version.CurrentCVEDescriptions {
					if cve.Name == "" && cve.Description == "" {
						handleErr(&service, fmt.Errorf("invalid CVE description"))
						return
					}
				}
				if version.CurrentCVEs == 0 {
					version.CurrentCVEs = len(version.CurrentCVEDescriptions)
				} else if version.CurrentCVEs < 0 {
					handleErr(&service, fmt.Errorf("invalid version"))
					return
				}
				for _, cve := range version.LatestCVEDescriptions {
					if cve.Name == "" && cve.Description == "" {
						handleErr(&service, fmt.Errorf("invalid CVE description"))
						return
					}
				}
				if version.LatestCVEs == 0 {
					version.LatestCVEs = len(version.LatestCVEDescriptions)
				} else if version.LatestCVEs < 0 {
					handleErr(&service, fmt.Errorf("invalid version"))
					return
				}
				versions = append(versions, version)

			case "actionGroup":
				var i ActionGroupInstruction
				if err = decodeServiceInstruction(result, &i); err != nil {
					handleErr(&service, fmt.Errorf("invalid instruction - %s", err))
					return
				}
				actionGroup := i.ServiceActionGroupState
				if actionGroup.Name == "" {
					handleErr(&service, fmt.Errorf("invalid action group"))
					return
				}
				if actionGroup.Time, err = decodeTime(i.Time); err != nil {
					handleErr(&service, err)
					return
				}
				if actionGroup.Order == "" {
					actionGroup.Order = actionGroup.Name
				}
				actionGroups = append(actionGroups, actionGroup)

			case "action":
				var i ActionInstruction
				if err = decodeServiceInstruction(result, &i); err != nil {
					handleErr(&service, fmt.Errorf("invalid instruction - %s", err))
					return
				}
				action := i.ServiceActionState
				if action.Name == "" {
					handleErr(&service, fmt.Errorf("invalid action"))
					return
				}
				if action.Time, err = decodeTime(i.Time); err != nil {
					handleErr(&service, err)
					return
				}
				if action.Order == "" {
					action.Order = action.Name
				}
				if (action.URL == "" && action.Script == "") || (action.URL != "" && action.Script != "") {
					handleErr(&service, fmt.Errorf("invalid action"))
					return
				}
				actions = append(actions, action)

			case "debug":
				// ignore

			default:
				handleErr(&service, fmt.Errorf("invalid instruction type \"%s\"", i.Type))
				return
			}
		}

		serviceState.lock.Lock()
		scriptState.error = nil
		// current uptime statuses
		scriptState.uptimeStatuses = make(map[string]ServiceUptimeStatusState, len(scriptState.uptimeStatuses))
		for _, uptimeStatus := range uptimeStatuses {
			if existing, ok := scriptState.uptimeStatuses[uptimeStatus.Name]; !ok || existing.Time.Before(uptimeStatus.Time) {
				scriptState.uptimeStatuses[uptimeStatus.Name] = uptimeStatus
			}
		}
		// current versions
		scriptState.versions = make(map[string]ServiceVersionState, len(scriptState.versions))
		for _, version := range versions {
			if existing, ok := scriptState.versions[version.Name]; !ok || existing.Time.Before(version.Time) {
				scriptState.versions[version.Name] = version
			}
		}
		// current action groups
		scriptState.actionGroups = make(map[string]ServiceActionGroupState, len(scriptState.actionGroups))
		for _, actionGroup := range actionGroups {
			if existing, ok := scriptState.actionGroups[actionGroup.Name]; !ok || existing.Time.Before(actionGroup.Time) {
				scriptState.actionGroups[actionGroup.Name] = actionGroup
			}
		}
		// current actions
		scriptState.actions = make(map[string]ServiceActionState, len(scriptState.actions))
		for _, action := range actions {
			if existing, ok := scriptState.actions[action.Name]; !ok || existing.Time.Before(action.Time) {
				scriptState.actions[action.Name] = action
			}
		}
		// uptime statuses to be alerted
		previousUptimeStatuses := make(map[string]ServiceUptimeStatusState, len(scriptState.uptimeStatuses))
		for name := range scriptState.uptimeStatuses {
			if existing, ok := serviceState.uptimeStatuses[name]; ok {
				previousUptimeStatuses[name] = existing
			}
		}
		for _, uptimeStatus := range uptimeStatuses {
			latest := scriptState.uptimeStatuses[uptimeStatus.Name]
			if existing, ok := previousUptimeStatuses[uptimeStatus.Name]; (!ok || existing.Time.Before(uptimeStatus.Time)) && latest.Time.After(uptimeStatus.Time) {
				previousUptimeStatuses[uptimeStatus.Name] = uptimeStatus
			}
		}
		newlyUnhealthyUptimeStatuses := make([]ServiceUptimeStatusState, 0, len(scriptState.uptimeStatuses))
		for name, uptimeStatus := range scriptState.uptimeStatuses {
			if existing, ok := previousUptimeStatuses[name]; uptimeStatus.Unhealthy() && (!ok || (existing.Time.Before(uptimeStatus.Time) && !existing.Unhealthy())) {
				newlyUnhealthyUptimeStatuses = append(newlyUnhealthyUptimeStatuses, uptimeStatus)
			}
		}
		// versions to be alerted
		previousVersions := make(map[string]ServiceVersionState, len(scriptState.versions))
		for name := range scriptState.versions {
			if existing, ok := serviceState.versions[name]; ok {
				previousVersions[name] = existing
			}
		}
		for _, version := range versions {
			latest := scriptState.versions[version.Name]
			if existing, ok := previousVersions[version.Name]; (!ok || existing.Time.Before(version.Time)) && latest.Time.After(version.Time) {
				previousVersions[version.Name] = version
			}
		}
		newlyOutdatedVersions := make([]ServiceVersionState, 0, len(scriptState.versions))
		newlyVulnerableVersions := make([]ServiceVersionState, 0, len(scriptState.versions))
		for name, version := range scriptState.versions {
			existing, ok := previousVersions[name]
			if version.Outdated() && (!ok || (existing.Time.Before(version.Time) && !existing.Outdated())) {
				newlyOutdatedVersions = append(newlyOutdatedVersions, version)
			}
			if version.Vulnerable() && (!ok || (existing.Time.Before(version.Time) && !existing.Vulnerable())) {
				newlyVulnerableVersions = append(newlyVulnerableVersions, version)
			}
		}
		updateServiceState(serviceState)
		serviceState.lock.Unlock()

		if len(newlyUnhealthyUptimeStatuses) > 0 {
			alertUnhealthyUptimeStatuses(s.app, service, script, newlyUnhealthyUptimeStatuses)
		}
		if len(newlyOutdatedVersions) > 0 {
			alertOutdatedVersions(s.app, service, script, newlyOutdatedVersions)
		}
		if len(newlyVulnerableVersions) > 0 {
			alertVulnerableVersions(s.app, service, script, newlyVulnerableVersions)
		}
	}
}

func updateServiceState(serviceState *serviceState) {
	serviceState.uptimeStatuses = make(map[string]ServiceUptimeStatusState, len(serviceState.uptimeStatuses))
	serviceState.versions = make(map[string]ServiceVersionState, len(serviceState.versions))
	serviceState.actionGroups = make(map[string]ServiceActionGroupState, len(serviceState.actionGroups))
	serviceState.actions = make(map[string]ServiceActionState, len(serviceState.actions))
	for _, scriptState := range serviceState.scripts {
		for name, uptimeStatus := range scriptState.uptimeStatuses {
			if existing, ok := serviceState.uptimeStatuses[name]; !ok || existing.Time.Before(uptimeStatus.Time) {
				serviceState.uptimeStatuses[name] = uptimeStatus
			}
		}
		for name, version := range scriptState.versions {
			if existing, ok := serviceState.versions[name]; !ok || existing.Time.Before(version.Time) {
				serviceState.versions[name] = version
			}
		}
		for name, actionGroup := range scriptState.actionGroups {
			if existing, ok := serviceState.actionGroups[name]; !ok || existing.Time.Before(actionGroup.Time) {
				serviceState.actionGroups[name] = actionGroup
			}
		}
		for name, action := range scriptState.actions {
			if existing, ok := serviceState.actions[name]; !ok || existing.Time.Before(action.Time) {
				serviceState.actions[name] = action
			}
		}
	}
	serviceState.uptimeStatusesSorted = make([]ServiceUptimeStatusState, 0, len(serviceState.uptimeStatuses))
	for _, uptimeStatus := range serviceState.uptimeStatuses {
		serviceState.uptimeStatusesSorted = append(serviceState.uptimeStatusesSorted, uptimeStatus)
	}
	sort.Sort(ByUptimeStatusOrder(serviceState.uptimeStatusesSorted))
	serviceState.versionsSorted = make([]ServiceVersionState, 0, len(serviceState.versions))
	for _, version := range serviceState.versions {
		serviceState.versionsSorted = append(serviceState.versionsSorted, version)
	}
	sort.Sort(ByVersionOrder(serviceState.versionsSorted))
	serviceState.actionGroupsSorted = make([]ServiceActionGroupState, 0, len(serviceState.actionGroups))
	for _, actionGroup := range serviceState.actionGroups {
		serviceState.actionGroupsSorted = append(serviceState.actionGroupsSorted, actionGroup)
	}
	sort.Sort(ByActionGroupOrder(serviceState.actionGroupsSorted))
	serviceState.actionsSorted = make([]ServiceActionState, 0, len(serviceState.actions))
	for _, action := range serviceState.actions {
		serviceState.actionsSorted = append(serviceState.actionsSorted, action)
	}
	sort.Sort(ByActionOrder(serviceState.actionsSorted))
}
