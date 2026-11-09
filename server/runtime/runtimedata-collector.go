package runtime

import (
	"fmt"

	"github.com/2manyvcos/paranal/server/collector"
	"github.com/2manyvcos/paranal/server/dynamicconfig"
	gojpl "github.com/jplorg/jpl/go/v2"
)

func (d *RuntimeData) createCollectorFromDynamicConfig(id string, config dynamicconfig.CollectorConfig) (*collector.Collector, error) {
	program, jplErr := gojpl.Parse(config.Source, nil)
	if jplErr != nil {
		return nil, jplErr
	}

	triggers := make([]collector.Trigger, len(config.Triggers))
	for i, t := range config.Triggers {
		switch t.Trigger {
		case "start":
			triggers[i] = collector.OnStart()

		case "schedule":
			triggers[i] = collector.OnSchedule(d.cronScheduler, t.Crontab)

		case "change":
			triggers[i] = collector.OnChange(d.eventBroker, t.PathPattern)

		default:
			return nil, fmt.Errorf("unsupported trigger %q", t.Trigger)
		}
	}

	return collector.NewContext(d.ctx, d.cfg, d.store, id, program, triggers...)
}
