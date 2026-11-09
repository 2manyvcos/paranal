package dynamicconfig

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-viper/mapstructure/v2"
	"golang.org/x/mod/semver"
)

const collectorConfigV1Version = "v1.0"

type collectorConfigV1 struct {
	ID       string            `mapstructure:"id"`
	Triggers []triggerConfigV1 `mapstructure:"triggers"`
	Src      string            `mapstructure:"src"`
}

type triggerConfigV1 struct {
	Trigger string         `mapstructure:"trigger"`
	Params  map[string]any `mapstructure:",remain"`
}

type startTriggerParamsV1 struct{}

type scheduleTriggerParamsV1 struct {
	Cron string `mapstructure:"cron"`
}

type changeTriggerParamsV1 struct {
	Pattern string `mapstructure:"pattern"`
}

func parseCollectorConfig(version string, config any) (Config, error) {
	switch semver.Major(version) {
	case "v1":
		if semver.Compare(version, collectorConfigV1Version) > 0 {
			return nil, fmt.Errorf("unsupported version %q", version)
		}
		return parseCollectorConfigV1(config)

	default:
		return nil, fmt.Errorf("unsupported version %q", version)
	}
}

func parseCollectorConfigV1(config any) (Config, error) {
	var c collectorConfigV1
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &c})
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	var result CollectorConfig

	if c.ID == "" {
		return nil, fmt.Errorf("missing \"id\"")
	}
	result.ID = c.ID

	result.Triggers = make([]TriggerConfig, len(c.Triggers))
	for i, trigger := range c.Triggers {
		result.Triggers[i], err = parseTriggerConfigV1(trigger)
		if err != nil {
			return nil, fmt.Errorf("parsing trigger: %w", err)
		}
	}

	if c.Src == "" {
		return nil, fmt.Errorf("missing \"src\"")
	}
	result.Source = c.Src

	return result, nil
}

func parseTriggerConfigV1(config triggerConfigV1) (TriggerConfig, error) {
	switch config.Trigger {
	case "start":
		trigger, err := parseStartTriggerV1(config.Params)
		if err != nil {
			return TriggerConfig{}, err
		}
		return trigger, nil

	case "schedule":
		trigger, err := parseScheduleTriggerV1(config.Params)
		if err != nil {
			return TriggerConfig{}, err
		}
		return trigger, nil

	case "change":
		trigger, err := parseChangeTriggerV1(config.Params)
		if err != nil {
			return TriggerConfig{}, err
		}
		return trigger, nil

	default:
		return TriggerConfig{}, fmt.Errorf("unsupported trigger %q", config.Trigger)
	}
}

func parseStartTriggerV1(config any) (TriggerConfig, error) {
	var p startTriggerParamsV1
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &p})
	if err != nil {
		return TriggerConfig{}, err
	}
	err = decoder.Decode(config)
	if err != nil {
		return TriggerConfig{}, err
	}

	var result TriggerConfig
	result.Trigger = "start"

	return result, nil
}

func parseScheduleTriggerV1(config any) (TriggerConfig, error) {
	var p scheduleTriggerParamsV1
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &p})
	if err != nil {
		return TriggerConfig{}, err
	}
	err = decoder.Decode(config)
	if err != nil {
		return TriggerConfig{}, err
	}

	var result TriggerConfig
	result.Trigger = "schedule"

	if err := gocron.NewDefaultCron(true).IsValid(p.Cron, time.Local, time.Now()); err != nil {
		return TriggerConfig{}, fmt.Errorf("parsing cron schedule: %w", err)
	}
	result.Crontab = p.Cron

	return result, nil
}

func parseChangeTriggerV1(config any) (TriggerConfig, error) {
	var p changeTriggerParamsV1
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &p})
	if err != nil {
		return TriggerConfig{}, err
	}
	err = decoder.Decode(config)
	if err != nil {
		return TriggerConfig{}, err
	}

	var result TriggerConfig
	result.Trigger = "change"

	if p.Pattern == "" {
		return TriggerConfig{}, fmt.Errorf("missing \"pattern\"")
	}
	result.PathPattern = p.Pattern

	return result, nil
}
