package state

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/2manyvcos/paranal/server/scripts"
	"github.com/go-viper/mapstructure/v2"
)

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
				if version.CurrentVersion == "" || version.LatestVersion == "" {
					handleErr(fmt.Errorf("invalid version"))
					return
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
		uptimeStatusesToBeAlerted := make([]ServiceUptimeStatusState, 0, len(scriptState.uptimeStatuses))
		for name, uptimeStatus := range scriptState.uptimeStatuses {
			if existing, ok := previousUptimeStatuses[name]; uptimeStatus.Status != ServiceUptimeStatusUp && (!ok || (existing.Time.Before(uptimeStatus.Time) && existing.Status == ServiceUptimeStatusUp)) {
				uptimeStatusesToBeAlerted = append(uptimeStatusesToBeAlerted, uptimeStatus)
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
		versionsToBeAlerted := make([]ServiceVersionState, 0, len(scriptState.versions))
		for name, version := range scriptState.versions {
			if existing, ok := previousVersions[name]; version.CurrentVersion != version.LatestVersion && (!ok || (existing.Time.Before(version.Time) && existing.CurrentVersion == existing.LatestVersion)) {
				versionsToBeAlerted = append(versionsToBeAlerted, version)
			}
		}
		updateServiceState(serviceState)
		serviceState.lock.Unlock()

		if len(uptimeStatusesToBeAlerted) > 0 {
			alertUptimeStatuses(s.app, script, uptimeStatusesToBeAlerted)
		}
		if len(versionsToBeAlerted) > 0 {
			alertVersions(s.app, script, versionsToBeAlerted)
		}
	}
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

type ByUptimeStatusOrder []ServiceUptimeStatusState

func (a ByUptimeStatusOrder) Len() int           { return len(a) }
func (a ByUptimeStatusOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByUptimeStatusOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

type VersionInstruction struct {
	Type string `mapstructure:"type"`
	ServiceVersionState
	Time any `mapstructure:"time"`
}

type ByVersionOrder []ServiceVersionState

func (a ByVersionOrder) Len() int           { return len(a) }
func (a ByVersionOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVersionOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

type ContextSectionInstruction struct {
	Type string `mapstructure:"type"`
	ServiceContextSectionState
	Time any `mapstructure:"time"`
}

type ByContextSectionOrder []ServiceContextSectionState

func (a ByContextSectionOrder) Len() int           { return len(a) }
func (a ByContextSectionOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByContextSectionOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

type ContextOptionInstruction struct {
	Type string `mapstructure:"type"`
	ServiceContextOptionState
	Time any `mapstructure:"time"`
}

type ByContextOptionOrder []ServiceContextOptionState

func (a ByContextOptionOrder) Len() int           { return len(a) }
func (a ByContextOptionOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByContextOptionOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

func decodeServiceInstruction(input, output any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: output, Squash: true, IgnoreUntaggedFields: true})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func decodeTime(input any) (time.Time, error) {
	if input == nil {
		return time.Now(), nil
	}
	switch v := input.(type) {
	case string:
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC3339, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC1123Z, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC1123, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC850, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC822Z, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.RFC822, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.DateTime, v); err == nil {
			return t, nil
		} else if t, err = time.Parse(time.DateOnly, v); err == nil {
			return t, nil
		} else {
			return time.Time{}, fmt.Errorf("invalid time")
		}

	case float64:
		return time.Unix(int64(v), 0), nil

	default:
		return time.Time{}, fmt.Errorf("invalid time")
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
