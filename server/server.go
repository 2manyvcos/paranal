package server

import (
	"context"
	"log"

	meta "github.com/2manyvcos/paranal"
	"github.com/2manyvcos/paranal/server/config"
	"github.com/2manyvcos/paranal/server/events"
	"github.com/2manyvcos/paranal/server/storage"
	"github.com/go-co-op/gocron/v2"
)

func Run(ctx context.Context) error {
	log.Printf("welcome to %v (version %v)\n", meta.Meta.Name, meta.Meta.Version)

	cfg, err := config.Load()
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

	storage, err := storage.New(cfg, eventBroker)
	if err != nil {
		return err
	}
	defer storage.Shutdown()

	// TODO: add additional setup code here
	// TODO: use collector.NewContext(ctx, …) so that collectors are stopped early on signal instead of waiting for the HTTP server to exit

	cronScheduler.Start()

	return Serve(ctx, cfg)
}
