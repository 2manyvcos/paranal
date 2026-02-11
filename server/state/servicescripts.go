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

func newServiceScriptRunner(s *State, serviceState *serviceState, scriptState *serviceScriptState, script schema.ServiceScript) func() {
	handleErr := func(err error) {
		log.Printf("Error running script \"%s\" for service \"%s\" - %s", script.ServiceID, script.ID, err)
		serviceState.lock.Lock()
		scriptState.error = err
		scriptState.uptimeStatuses = nil
		scriptState.versions = nil
		scriptState.contextSections = nil
		scriptState.contextOptions = nil
		updateServiceState(serviceState)
		serviceState.lock.Unlock()
		alertScriptError(s.app, script, err)
	}

	return func() {
		results, err := scripts.RunServiceScript(s.app, script.ServiceID, script.Source)
		if err != nil {
			handleErr(err)
			return
		}

		var uptimeStatuses []ServiceUptimeStatusState
		var versions []ServiceVersionState
		var contextSections []ServiceContextSectionState
		var contextOptions []ServiceContextOptionState

		for _, result := range results {
			var i GenericInstruction
			err := mapstructure.Decode(result, &i)
			if err != nil {
				handleErr(err)
				return
			}

			switch i.Type {
			case "uptimeStatus":
				var i UptimeStatusInstruction
				if err = decodeServiceInstruction(result, &i); err != nil {
					handleErr(fmt.Errorf("invalid instruction - %s", err))
					return
				}
				uptimeStatus := i.ServiceUptimeStatusState
				if uptimeStatus.Time, err = decodeTime(i.Time); err != nil {
					handleErr(err)
					return
				}
				if uptimeStatus.Order == "" {
					uptimeStatus.Order = uptimeStatus.Name
				}
				if status, ok := ServiceUptimeStatusCodes[i.Status]; ok {
					uptimeStatus.Status = status
				} else {
					handleErr(fmt.Errorf("invalid status"))
					return
				}
				uptimeStatuses = append(uptimeStatuses, uptimeStatus)

			case "version":
				var i VersionInstruction
				if err = decodeServiceInstruction(result, &i); err != nil {
					handleErr(fmt.Errorf("invalid instruction - %s", err))
					return
				}
				version := i.ServiceVersionState
				if version.Time, err = decodeTime(i.Time); err != nil {
					handleErr(err)
					return
				}
				if version.Order == "" {
					version.Order = version.Name
				}
				if version.CurrentVersion == "" {
					handleErr(fmt.Errorf("invalid version"))
					return
				}
				if i.Status == "" {
					if version.LatestVersion == "" {
						handleErr(fmt.Errorf("invalid version"))
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
					handleErr(fmt.Errorf("invalid status"))
					return
				}
				for _, cve := range version.CurrentVersionCVEs {
					if cve.Name == "" && cve.Description == "" {
						handleErr(fmt.Errorf("invalid CVE"))
						return
					}
				}
				for _, cve := range version.LatestVersionCVEs {
					if cve.Name == "" && cve.Description == "" {
						handleErr(fmt.Errorf("invalid CVE"))
						return
					}
				}
				versions = append(versions, version)

			case "contextSection":
				var i ContextSectionInstruction
				if err = decodeServiceInstruction(result, &i); err != nil {
					handleErr(fmt.Errorf("invalid instruction - %s", err))
					return
				}
				contextSection := i.ServiceContextSectionState
				if contextSection.Time, err = decodeTime(i.Time); err != nil {
					handleErr(err)
					return
				}
				if contextSection.Order == "" {
					contextSection.Order = contextSection.Name
				}
				contextSections = append(contextSections, contextSection)

			case "contextOption":
				var i ContextOptionInstruction
				if err = decodeServiceInstruction(result, &i); err != nil {
					handleErr(fmt.Errorf("invalid instruction - %s", err))
					return
				}
				contextOption := i.ServiceContextOptionState
				if contextOption.Time, err = decodeTime(i.Time); err != nil {
					handleErr(err)
					return
				}
				if contextOption.Order == "" {
					contextOption.Order = contextOption.Name
				}
				if (contextOption.URL == "" && contextOption.Script == "") || (contextOption.URL != "" && contextOption.Script != "") {
					handleErr(fmt.Errorf("invalid context option"))
					return
				}
				contextOptions = append(contextOptions, contextOption)

			case "debug":
				// ignore

			default:
				handleErr(fmt.Errorf("invalid instruction type \"%s\"", i.Type))
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
		// current context sections
		scriptState.contextSections = make(map[string]ServiceContextSectionState, len(scriptState.contextSections))
		for _, contextSection := range contextSections {
			if existing, ok := scriptState.contextSections[contextSection.Name]; !ok || existing.Time.Before(contextSection.Time) {
				scriptState.contextSections[contextSection.Name] = contextSection
			}
		}
		// current context options
		scriptState.contextOptions = make(map[string]ServiceContextOptionState, len(scriptState.contextOptions))
		for _, contextOption := range contextOptions {
			if existing, ok := scriptState.contextOptions[contextOption.Name]; !ok || existing.Time.Before(contextOption.Time) {
				scriptState.contextOptions[contextOption.Name] = contextOption
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
		newlyDownUptimeStatuses := make([]ServiceUptimeStatusState, 0, len(scriptState.uptimeStatuses))
		for name, uptimeStatus := range scriptState.uptimeStatuses {
			if existing, ok := previousUptimeStatuses[name]; !uptimeStatus.Up() && (!ok || (existing.Time.Before(uptimeStatus.Time) && existing.Up())) {
				newlyDownUptimeStatuses = append(newlyDownUptimeStatuses, uptimeStatus)
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
			if !version.UpToDate() && (!ok || (existing.Time.Before(version.Time) && existing.UpToDate())) {
				newlyOutdatedVersions = append(newlyOutdatedVersions, version)
			}
			if version.Vulnerable() && (!ok || (existing.Time.Before(version.Time) && !existing.Vulnerable())) {
				newlyVulnerableVersions = append(newlyVulnerableVersions, version)
			}
		}
		updateServiceState(serviceState)
		serviceState.lock.Unlock()

		if len(newlyDownUptimeStatuses) > 0 {
			alertDownUptimeStatuses(s.app, script, newlyDownUptimeStatuses)
		}
		if len(newlyOutdatedVersions) > 0 {
			alertOutdatedVersions(s.app, script, newlyOutdatedVersions)
		}
		if len(newlyVulnerableVersions) > 0 {
			alertVulnerableVersions(s.app, script, newlyVulnerableVersions)
		}
	}
}

func updateServiceState(serviceState *serviceState) {
	serviceState.uptimeStatuses = make(map[string]ServiceUptimeStatusState, len(serviceState.uptimeStatuses))
	serviceState.versions = make(map[string]ServiceVersionState, len(serviceState.versions))
	serviceState.contextSections = make(map[string]ServiceContextSectionState, len(serviceState.contextSections))
	serviceState.contextOptions = make(map[string]ServiceContextOptionState, len(serviceState.contextOptions))
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
		for name, contextSection := range scriptState.contextSections {
			if existing, ok := serviceState.contextSections[name]; !ok || existing.Time.Before(contextSection.Time) {
				serviceState.contextSections[name] = contextSection
			}
		}
		for name, contextOption := range scriptState.contextOptions {
			if existing, ok := serviceState.contextOptions[name]; !ok || existing.Time.Before(contextOption.Time) {
				serviceState.contextOptions[name] = contextOption
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
	serviceState.contextSectionsSorted = make([]ServiceContextSectionState, 0, len(serviceState.contextSections))
	for _, contextSection := range serviceState.contextSections {
		serviceState.contextSectionsSorted = append(serviceState.contextSectionsSorted, contextSection)
	}
	sort.Sort(ByContextSectionOrder(serviceState.contextSectionsSorted))
	serviceState.contextOptionsSorted = make([]ServiceContextOptionState, 0, len(serviceState.contextOptions))
	for _, contextOption := range serviceState.contextOptions {
		serviceState.contextOptionsSorted = append(serviceState.contextOptionsSorted, contextOption)
	}
	sort.Sort(ByContextOptionOrder(serviceState.contextOptionsSorted))
}
