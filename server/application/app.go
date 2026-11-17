package application

import (
	"context"
	"fmt"
	"log"
	"maps"
	"net/url"
	"slices"
	"sync"

	"github.com/2manyvcos/paranal/server/collector"
	"github.com/2manyvcos/paranal/server/dynamicconfig"
	"github.com/2manyvcos/paranal/server/events"
	"github.com/2manyvcos/paranal/server/staticconfig"
	"github.com/2manyvcos/paranal/server/storage"
	dashboard_spec "github.com/2manyvcos/paranal/spec/dashboard"
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
	dashboards    map[string]dashboard_spec.Dashboard
}

func (a *App) setup(dynamicConfigs []dynamicconfig.Config) error {
	newCollectors := make(map[string]*collector.Collector, len(dynamicConfigs))

	for _, config := range dynamicConfigs {
		switch config.Type() {
		case "collector":
			id := config.ID()
			if c, ok := a.collectors[id]; ok {
				if c != nil {
					c.Shutdown()
				}
				delete(a.collectors, id)
			}
			delete(newCollectors, id)
			c, err := a.createCollectorFromDynamicConfig(id, config.(dynamicconfig.CollectorConfig))
			if err != nil {
				log.Printf("setting up collector %q: %v\n", id, err)
				a.collectors[id] = nil
				continue
			}
			a.collectors[id] = c
			newCollectors[id] = c

		case "dashboard":
			id := config.ID()
			a.dashboards[id] = config.(dynamicconfig.DashboardConfig).Dashboard

		default:
			return fmt.Errorf("unsupported type %q", config.Type())
		}
	}

	_, err := a.Store.CleanupCollectorStates(slices.Collect(maps.Keys(a.collectors)))
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
	changedDashboards := make(map[string]struct{}, len(changes))

	for _, change := range changes {
		switch change.Config.Type() {
		case "collector":
			id := change.ID()
			if c, ok := a.collectors[id]; ok {
				if c != nil {
					c.Shutdown()
				}
				delete(a.collectors, id)
			}
			delete(removedCollectors, id)
			delete(newCollectors, id)
			if change.Removed {
				removedCollectors[id] = struct{}{}
				continue
			}
			c, err := a.createCollectorFromDynamicConfig(id, change.Config.(dynamicconfig.CollectorConfig))
			if err != nil {
				log.Printf("applying config change: setting up collector %q: %v\n", id, err)
				a.collectors[id] = nil
				continue
			}
			a.collectors[id] = c
			newCollectors[id] = c

		case "dashboard":
			id := change.ID()
			changedDashboards[id] = struct{}{}
			if change.Removed {
				delete(a.dashboards, id)
				continue
			}
			a.dashboards[id] = change.Config.(dynamicconfig.DashboardConfig).Dashboard

		default:
			log.Printf("applying config change: unsupported type %q", change.Config.Type())
		}
	}

	var eventBundle []events.Event

	for id := range removedCollectors {
		events, err := a.Store.PutCollectorStates(id, nil)
		if err != nil {
			log.Printf("cleaning up states: %v", err)
		}

		eventBundle = append(eventBundle, events...)
	}

	eventBundle = slices.Grow(eventBundle, len(changedDashboards))
	for id := range changedDashboards {
		eventBundle = append(eventBundle, events.NewChangeEvent(fmt.Sprintf("v1/dashboards/%s", url.PathEscape(id))))
	}

	a.EventBroker.Publish(eventBundle...)

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
	a.dashboards = make(map[string]dashboard_spec.Dashboard)

	err := a.setup(dynamicConfigs)
	if err != nil {
		return nil, err
	}

	go a.watch(dynamicConfigChanges)

	return &a, nil
}
