package application

import (
	"context"
	"fmt"
	"log"
	"maps"
	"slices"
	"sync"

	"github.com/2manyvcos/paranal/server/collector"
	"github.com/2manyvcos/paranal/server/dynamicconfig"
	"github.com/2manyvcos/paranal/server/events"
	"github.com/2manyvcos/paranal/server/staticconfig"
	"github.com/2manyvcos/paranal/server/storage"
	"github.com/go-co-op/gocron/v2"
)

type App struct {
	Config        *staticconfig.Config
	CronScheduler gocron.Scheduler
	EventBroker   *events.EventBroker
	Store         *storage.Storage
	ctx           context.Context
	lock          sync.Mutex
	collectors    map[string]*collector.Collector
}

func (a *App) setup(dynamicConfigs []dynamicconfig.Config) error {
	newCollectors := make(map[string]*collector.Collector, len(dynamicConfigs))

	for _, config := range dynamicConfigs {
		switch config.Kind() {
		case "collector":
			collectorConfig := config.(dynamicconfig.CollectorConfig)
			if c, ok := a.collectors[collectorConfig.ID]; ok {
				if c != nil {
					c.Shutdown()
				}
				delete(a.collectors, collectorConfig.ID)
			}
			delete(newCollectors, collectorConfig.ID)
			c, err := a.createCollectorFromDynamicConfig(collectorConfig.ID, collectorConfig)
			if err != nil {
				log.Printf("setting up collector %q: %v\n", collectorConfig.ID, err)
				a.collectors[collectorConfig.ID] = nil
				continue
			}
			a.collectors[collectorConfig.ID] = c
			newCollectors[collectorConfig.ID] = c

		default:
			return fmt.Errorf("unsupported kind %q", config.Kind())
		}
	}

	err := a.Store.CleanupCollectorStates(slices.Collect(maps.Keys(a.collectors)))
	if err != nil {
		return fmt.Errorf("cleaning up states: %w", err)
	}

	for _, c := range newCollectors {
		c.Start()
	}

	return nil
}

func (a *App) applyChanges(changes []dynamicconfig.ConfigChange) {
	a.lock.Lock()
	defer a.lock.Unlock()

	removedCollectors := make(map[string]struct{}, len(changes))
	newCollectors := make(map[string]*collector.Collector, len(changes))

	for _, change := range changes {
		switch change.Config.Kind() {
		case "collector":
			collectorConfig := change.Config.(dynamicconfig.CollectorConfig)
			if c, ok := a.collectors[collectorConfig.ID]; ok {
				if c != nil {
					c.Shutdown()
				}
				delete(a.collectors, collectorConfig.ID)
			}
			delete(removedCollectors, collectorConfig.ID)
			delete(newCollectors, collectorConfig.ID)
			if change.Removed {
				removedCollectors[collectorConfig.ID] = struct{}{}
				continue
			}
			c, err := a.createCollectorFromDynamicConfig(collectorConfig.ID, collectorConfig)
			if err != nil {
				log.Printf("applying config change: setting up collector %q: %v\n", collectorConfig.ID, err)
				a.collectors[collectorConfig.ID] = nil
				continue
			}
			a.collectors[collectorConfig.ID] = c
			newCollectors[collectorConfig.ID] = c

		default:
			log.Printf("applying config change: unsupported kind %q", change.Config.Kind())
		}
	}

	for id := range removedCollectors {
		err := a.Store.PutCollectorStates(id, nil)
		if err != nil {
			log.Printf("cleaning up states: %v", err)
		}
	}

	for _, c := range newCollectors {
		c.Start()
	}
}

func (a *App) watch(dynamicConfigChanges <-chan []dynamicconfig.ConfigChange) {
	for {
		select {
		case <-a.ctx.Done():
			return

		case changes := <-dynamicConfigChanges:
			log.Println("applying config changes")
			a.applyChanges(changes)
		}
	}
}

func New(
	ctx context.Context,
	cfg *staticconfig.Config,
	cronScheduler gocron.Scheduler,
	eventBroker *events.EventBroker,
	store *storage.Storage,
	dynamicConfigs []dynamicconfig.Config,
	dynamicConfigChanges <-chan []dynamicconfig.ConfigChange,
) (*App, error) {
	var a App
	a.Config = cfg
	a.CronScheduler = cronScheduler
	a.EventBroker = eventBroker
	a.Store = store
	a.ctx = ctx
	a.collectors = make(map[string]*collector.Collector)

	err := a.setup(dynamicConfigs)
	if err != nil {
		return nil, err
	}

	go a.watch(dynamicConfigChanges)

	return &a, nil
}
