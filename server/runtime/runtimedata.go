package runtime

import (
	"context"
	"fmt"
	"log"
	"maps"
	"net/url"
	"slices"
	"sync"

	"github.com/2manyvcos/paranal/server/collector"
	"github.com/2manyvcos/paranal/server/config"
	"github.com/2manyvcos/paranal/server/dynamicconfig"
	"github.com/2manyvcos/paranal/server/events"
	"github.com/2manyvcos/paranal/server/storage"
	"github.com/go-co-op/gocron/v2"
)

type RuntimeData struct {
	ctx           context.Context
	cfg           *config.Config
	cronScheduler gocron.Scheduler
	eventBroker   *events.EventBroker
	store         *storage.Storage
	lock          sync.Mutex
	collectors    map[string]*collector.Collector
}

func (d *RuntimeData) setup(dynamicConfigs []dynamicconfig.Config) error {
	newCollectors := make(map[string]*collector.Collector, len(dynamicConfigs))

	for _, config := range dynamicConfigs {
		switch config.Config() {
		case "collector":
			collectorConfig := config.(dynamicconfig.CollectorConfig)
			id := url.QueryEscape(collectorConfig.ID) + "@file"
			if c, ok := d.collectors[id]; ok {
				if c != nil {
					c.Shutdown()
				}
				delete(d.collectors, id)
			}
			delete(newCollectors, id)
			c, err := d.createCollectorFromDynamicConfig(id, collectorConfig)
			if err != nil {
				log.Printf("setting up collector %q: %v\n", id, err)
				d.collectors[id] = nil
				continue
			}
			d.collectors[id] = c
			newCollectors[id] = c

		default:
			return fmt.Errorf("unsupported config %q", config.Config())
		}
	}

	err := d.store.CleanupCollectorStates(slices.Collect(maps.Keys(d.collectors)))
	if err != nil {
		return fmt.Errorf("cleaning up states: %w", err)
	}

	for _, c := range newCollectors {
		c.Start()
	}

	return nil
}

func (d *RuntimeData) applyChanges(changes []dynamicconfig.ConfigChange) {
	d.lock.Lock()
	defer d.lock.Unlock()

	removedCollectors := make(map[string]struct{}, len(changes))
	newCollectors := make(map[string]*collector.Collector, len(changes))

	for _, change := range changes {
		switch change.Config.Config() {
		case "collector":
			collectorConfig := change.Config.(dynamicconfig.CollectorConfig)
			id := url.QueryEscape(collectorConfig.ID) + "@file"
			if c, ok := d.collectors[id]; ok {
				if c != nil {
					c.Shutdown()
				}
				delete(d.collectors, id)
			}
			delete(removedCollectors, id)
			delete(newCollectors, id)
			if change.Removed {
				removedCollectors[id] = struct{}{}
				continue
			}
			c, err := d.createCollectorFromDynamicConfig(id, collectorConfig)
			if err != nil {
				log.Printf("applying config change: setting up collector %q: %v\n", id, err)
				d.collectors[id] = nil
				continue
			}
			d.collectors[id] = c
			newCollectors[id] = c

		default:
			log.Printf("applying config change: unsupported config %q", change.Config.Config())
		}
	}

	for id := range removedCollectors {
		err := d.store.PutCollectorStates(id, nil)
		if err != nil {
			log.Printf("cleaning up states: %v", err)
		}
	}

	for _, c := range newCollectors {
		c.Start()
	}
}

func (d *RuntimeData) watch(dynamicConfigChanges <-chan []dynamicconfig.ConfigChange) {
	for {
		select {
		case <-d.ctx.Done():
			return

		case changes := <-dynamicConfigChanges:
			log.Println("applying config changes")
			d.applyChanges(changes)
		}
	}
}

func NewRuntimeData(
	ctx context.Context,
	cfg *config.Config,
	cronScheduler gocron.Scheduler,
	eventBroker *events.EventBroker,
	store *storage.Storage,
	dynamicConfigs []dynamicconfig.Config,
	dynamicConfigChanges <-chan []dynamicconfig.ConfigChange,
) (*RuntimeData, error) {
	var d RuntimeData
	d.ctx = ctx
	d.cfg = cfg
	d.cronScheduler = cronScheduler
	d.eventBroker = eventBroker
	d.store = store
	d.collectors = make(map[string]*collector.Collector)

	err := d.setup(dynamicConfigs)
	if err != nil {
		return nil, err
	}

	go d.watch(dynamicConfigChanges)

	return &d, nil
}
