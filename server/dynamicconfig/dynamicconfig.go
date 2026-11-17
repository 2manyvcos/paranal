package dynamicconfig

import (
	collector_spec "github.com/2manyvcos/paranal/spec/collector"
	dashboard_spec "github.com/2manyvcos/paranal/spec/dashboard"
)

type ConfigChange struct {
	Config
	Removed bool
}

type Config interface {
	Type() string
	ID() string
}

var _ Config = (*CollectorConfig)(nil)

type CollectorConfig struct{ collector_spec.Collector }

func (CollectorConfig) Type() string {
	return "collector"
}

func (c CollectorConfig) ID() string {
	return c.Collector.ID
}

var _ Config = (*DashboardConfig)(nil)

type DashboardConfig struct{ dashboard_spec.Dashboard }

func (DashboardConfig) Type() string {
	return "dashboard"
}

func (c DashboardConfig) ID() string {
	return c.Dashboard.ID
}
