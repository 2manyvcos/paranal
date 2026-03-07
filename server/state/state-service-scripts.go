package state

import (
	"log"

	"github.com/2manyvcos/paranal/server/schema"
)

func (s *State) ListServiceScriptStates(serviceID string) map[string]schema.ServiceScriptState {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return nil
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	result := make(map[string]schema.ServiceScriptState, len(service.scripts))
	for scriptID, script := range service.scripts {
		state := schema.ServiceScriptState{
			Running: script.running,
			Error:   script.error,
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

func (s *State) GetServiceScriptState(serviceID string, scriptID string) *schema.ServiceScriptState {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return nil
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	script, ok := service.scripts[scriptID]
	if !ok {
		return nil
	}
	state := schema.ServiceScriptState{
		Running: script.running,
		Error:   script.error,
	}
	if script.job != nil {
		if lastRun, err := script.job.LastRun(); err != nil && !lastRun.IsZero() {
			state.LastRun = &lastRun
		}
		if nextRun, err := script.job.NextRun(); err != nil && !nextRun.IsZero() {
			state.NextRun = &nextRun
		}
	}
	return &state
}

func (s *State) RunServiceScripts() {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	for serviceID, service := range s.services {
		service.lock.RLock()
		for scriptID, script := range service.scripts {
			if script.job != nil {
				err := script.job.RunNow()
				if err != nil {
					log.Printf("Error running script \"%s\" for service \"%s\" - %s\n", scriptID, serviceID, err)
				}
			}
		}
		service.lock.RUnlock()
	}
}

func (s *State) RunServiceScript(serviceID string, scriptID string) error {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return schema.ErrNotFound
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	script, ok := service.scripts[scriptID]
	if !ok {
		return schema.ErrNotFound
	}
	if script.running {
		return schema.ErrAlreadyRunning
	}
	if script.job != nil {
		err := script.job.RunNow()
		if err != nil {
			log.Printf("Error running script \"%s\" for service \"%s\" - %s\n", scriptID, serviceID, err)
		}
		return err
	}
	return nil
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
	updateServiceState(s, service, serviceID)
	if len(service.scripts) == 0 {
		delete(s.services, serviceID)
	}
}

func (s *State) ListServiceActionGroups(serviceID string) []schema.ServiceActionGroup {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return nil
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	result := make([]schema.ServiceActionGroup, len(service.actionGroupsSorted))
	for i, actionGroup := range service.actionGroupsSorted {
		result[i] = schema.ServiceActionGroup{
			Name: actionGroup.Name,
		}
	}
	return result
}

func (s *State) ListServiceActions(serviceID string) []schema.ServiceAction {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return nil
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	result := make([]schema.ServiceAction, len(service.actionsSorted))
	for i, action := range service.actionsSorted {
		result[i] = schema.ServiceAction{
			Name:             action.Name,
			URL:              action.URL,
			Script:           action.Script,
			Group:            action.Group,
			RestrictToAdmins: action.RestrictToAdmins,
		}
	}
	return result
}

func (s *State) GetServiceAction(serviceID string, actionName string) *schema.ServiceAction {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return nil
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	action, ok := service.actions[actionName]
	if !ok {
		return nil
	}
	result := schema.ServiceAction{
		Name:             action.Name,
		URL:              action.URL,
		Script:           action.Script,
		Group:            action.Group,
		RestrictToAdmins: action.RestrictToAdmins,
	}
	return &result
}

func (s *State) ListServiceHealthStatuses(serviceID string) []schema.ServiceHealthStatus {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return nil
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	result := make([]schema.ServiceHealthStatus, len(service.healthStatusesSorted))
	for i, healthStatus := range service.healthStatusesSorted {
		result[i] = schema.ServiceHealthStatus{
			Name:      healthStatus.Name,
			Status:    healthStatus.Status,
			Unhealthy: healthStatus.Unhealthy(),
		}
	}
	return result
}

func (s *State) ListServiceVersions(serviceID string) []schema.ServiceVersion {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return nil
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	result := make([]schema.ServiceVersion, len(service.versionsSorted))
	for i, version := range service.versionsSorted {
		result[i] = schema.ServiceVersion{
			Name:           version.Name,
			CurrentVersion: version.CurrentVersion,
			CurrentCVEs:    version.CurrentCVEs,
			LatestVersion:  version.LatestVersion,
			LatestCVEs:     version.LatestCVEs,
			Status:         version.Status,
			Outdated:       version.Outdated(),
			Vulnerable:     version.Vulnerable(),
		}
	}
	return result
}

func (s *State) GetServiceVersionDetails(serviceID string, versionName string) *schema.ServiceVersionDetails {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	service, ok := s.services[serviceID]
	if !ok {
		return nil
	}
	service.lock.RLock()
	defer service.lock.RUnlock()
	version, ok := service.versions[versionName]
	if !ok {
		return nil
	}
	currentCVEDescriptions := make([]schema.ServiceCVE, len(version.CurrentCVEDescriptions))
	for i, cveDescription := range version.CurrentCVEDescriptions {
		currentCVEDescriptions[i] = schema.ServiceCVE{
			Name:        cveDescription.Name,
			Description: cveDescription.Description,
		}
	}
	latestCVEDescriptions := make([]schema.ServiceCVE, len(version.LatestCVEDescriptions))
	for i, cveDescription := range version.LatestCVEDescriptions {
		latestCVEDescriptions[i] = schema.ServiceCVE{
			Name:        cveDescription.Name,
			Description: cveDescription.Description,
		}
	}
	result := schema.ServiceVersionDetails{
		CurrentVersionNotes:    version.CurrentVersionNotes,
		CurrentCVEDescriptions: currentCVEDescriptions,
		LatestVersionNotes:     version.LatestVersionNotes,
		LatestCVEDescriptions:  latestCVEDescriptions,
	}
	return &result
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
		updateServiceState(s, service, serviceID)
		service.lock.Unlock()
		delete(s.services, serviceID)
	}
}
