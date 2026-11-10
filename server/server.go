package server

import (
	"context"
	"log"
	"time"

	meta "github.com/2manyvcos/paranal"
	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/dynamicconfig"
	"github.com/2manyvcos/paranal/server/events"
	"github.com/2manyvcos/paranal/server/staticconfig"
	"github.com/2manyvcos/paranal/server/storage"
	"github.com/go-co-op/gocron/v2"
)

func Run(ctx context.Context) error {
	log.Printf("welcome to %v (version %v)\n", meta.Meta.Name, meta.Meta.Version)

	cfg, err := staticconfig.Load()
	if err != nil {
		return err
	}

	cronScheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	defer cronScheduler.Shutdown()

	eventBroker := events.NewEventBroker()
	defer eventBroker.Shutdown()

	store, err := storage.New(cfg, eventBroker)
	if err != nil {
		return err
	}
	defer store.Shutdown()

	dynamicConfigChanges := make(chan []dynamicconfig.ConfigChange)
	configLoader, dynamicConfigs, err := dynamicconfig.NewLoader(cfg.ConfigDir, 100*time.Millisecond, dynamicConfigChanges)
	if err != nil {
		return err
	}
	defer configLoader.Shutdown()

	app, err := application.New(ctx, cfg, cronScheduler, eventBroker, store, dynamicConfigs, dynamicConfigChanges)
	if err != nil {
		return err
	}

	if cfg.WatchConfigDir {
		configLoader.Watch()
	}

	cronScheduler.Start()

	return Serve(ctx, app)
}
