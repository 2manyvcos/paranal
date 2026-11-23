package collector_spec_v1

import (
	"fmt"
	"net/url"
	"time"

	collector_spec "github.com/2manyvcos/paranal/spec/collector"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-viper/mapstructure/v2"
)

type collector struct {
	ID       string    `mapstructure:"id"`
	Triggers []trigger `mapstructure:"triggers"`
	Src      string    `mapstructure:"src"`
}

type trigger struct {
	Type    string         `mapstructure:"type"`
	Payload map[string]any `mapstructure:",remain"`
}

func (collectorParser) ParseCollector(provider string, input any) (collector_spec.Collector, error) {
	var c collector
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &c})
	if err != nil {
		return collector_spec.Collector{}, err
	}
	err = decoder.Decode(input)
	if err != nil {
		return collector_spec.Collector{}, err
	}

	var result collector_spec.Collector
	result.Version = Version

	if c.ID == "" {
		return collector_spec.Collector{}, fmt.Errorf("missing \"id\"")
	}
	result.ID = fmt.Sprintf("%v@%v", url.QueryEscape(c.ID), url.QueryEscape(provider))

	if len(c.Triggers) == 0 {
		return collector_spec.Collector{}, fmt.Errorf("missing \"triggers\"")
	}
	result.Triggers = make([]collector_spec.Trigger, len(c.Triggers))
	for i, trigger := range c.Triggers {
		result.Triggers[i], err = parseTrigger(trigger)
		if err != nil {
			return collector_spec.Collector{}, fmt.Errorf("parsing \"triggers\": %w", err)
		}
	}

	if c.Src == "" {
		return collector_spec.Collector{}, fmt.Errorf("missing \"src\"")
	}
	result.Source = c.Src

	return result, nil
}

func parseTrigger(input trigger) (collector_spec.Trigger, error) {
	if input.Type == "" {
		return collector_spec.Trigger{}, fmt.Errorf("missing \"type\"")
	}

	switch input.Type {
	case "start":
		return parseStartTrigger(input.Payload)

	case "schedule":
		return parseScheduleTrigger(input.Payload)

	case "change":
		return parseChangeTrigger(input.Payload)

	default:
		return collector_spec.Trigger{}, fmt.Errorf("unsupported type %q", input.Type)
	}
}

type startTrigger struct{}

func parseStartTrigger(input any) (collector_spec.Trigger, error) {
	var t startTrigger
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &t})
	if err != nil {
		return collector_spec.Trigger{}, err
	}
	err = decoder.Decode(input)
	if err != nil {
		return collector_spec.Trigger{}, err
	}

	var result collector_spec.Trigger
	result.Type = "start"

	return result, nil
}

type scheduleTrigger struct {
	Cron string `mapstructure:"cron"`
}

func parseScheduleTrigger(input any) (collector_spec.Trigger, error) {
	var t scheduleTrigger
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &t})
	if err != nil {
		return collector_spec.Trigger{}, err
	}
	err = decoder.Decode(input)
	if err != nil {
		return collector_spec.Trigger{}, err
	}

	var result collector_spec.Trigger
	result.Type = "schedule"

	if t.Cron == "" {
		return collector_spec.Trigger{}, fmt.Errorf("missing \"cron\"")
	}
	if err := gocron.NewDefaultCron(true).IsValid(t.Cron, time.Local, time.Now()); err != nil {
		return collector_spec.Trigger{}, fmt.Errorf("parsing cron schedule: %w", err)
	}
	result.Crontab = t.Cron

	return result, nil
}

type changeTrigger struct {
	Pattern string `mapstructure:"pattern"`
}

func parseChangeTrigger(input any) (collector_spec.Trigger, error) {
	var t changeTrigger
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &t})
	if err != nil {
		return collector_spec.Trigger{}, err
	}
	err = decoder.Decode(input)
	if err != nil {
		return collector_spec.Trigger{}, err
	}

	var result collector_spec.Trigger
	result.Type = "change"

	if t.Pattern == "" {
		return collector_spec.Trigger{}, fmt.Errorf("missing \"pattern\"")
	}
	result.PathPattern = t.Pattern

	return result, nil
}
