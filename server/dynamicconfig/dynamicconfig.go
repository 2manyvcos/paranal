package dynamicconfig

import collector_spec "github.com/2manyvcos/paranal/server/spec/collector"

type ConfigChange struct {
	Config  Config
	Removed bool
}

type Config interface {
	Kind() string
}

var _ Config = (*CollectorConfig)(nil)

type CollectorConfig collector_spec.Collector

func (CollectorConfig) Kind() string {
	return "collector"
}
