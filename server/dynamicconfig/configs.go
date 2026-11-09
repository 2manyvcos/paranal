package dynamicconfig

type ConfigChange struct {
	Config  Config
	Removed bool
}

type Config interface {
	Config() string
}

var _ Config = (*CollectorConfig)(nil)

type TriggerConfig struct {
	Trigger     string
	Crontab     string
	PathPattern string
}

type CollectorConfig struct {
	ID       string
	Triggers []TriggerConfig
	Source   string
}

func (c CollectorConfig) Config() string {
	return "collector"
}
