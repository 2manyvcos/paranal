package dynamicconfig

import (
	collector_spec "github.com/2manyvcos/paranal/spec/collector"
)

type ConfigChange struct {
	Config
	Removed bool
}

type Config interface {
	Kind() string
	ID() string
}

var _ Config = (*CollectorConfig)(nil)

type CollectorConfig struct{ collector_spec.Collector }

func (CollectorConfig) Kind() string {
	return "collector"
}

func (c CollectorConfig) ID() string {
	return c.Collector.ID
}
