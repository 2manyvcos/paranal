package collector

import (
	"fmt"
	"time"

	"github.com/2manyvcos/paranal/server/events"
	"github.com/go-co-op/gocron/v2"
)

type Trigger interface {
	setup(collector *Collector) error
}

var _ Trigger = (*startTrigger)(nil)

type startTrigger struct{}

func (t startTrigger) setup(collector *Collector) error {
	time := time.Now().UTC().Format(time.RFC3339Nano)

	collector.scheduler.Schedule(map[string]any{
		"type": "start",
		"time": time,
	})

	return nil
}

func OnStart() Trigger {
	return startTrigger{}
}

var _ Trigger = (*scheduleTrigger)(nil)

type scheduleTrigger struct {
	cronScheduler gocron.Scheduler
	crontab       string
}

func (t scheduleTrigger) setup(collector *Collector) error {
	job, err := t.cronScheduler.NewJob(
		gocron.CronJob(t.crontab, true),
		gocron.NewTask(func() {
			time := time.Now().UTC().Format(time.RFC3339Nano)

			collector.scheduler.Schedule(map[string]any{
				"type": "schedule",
				"time": time,
			})
		}),
	)
	if err != nil {
		return err
	}

	go func() {
		<-collector.Done()
		t.cronScheduler.RemoveJob(job.ID())
	}()

	return nil
}

func OnSchedule(cronScheduler gocron.Scheduler, crontab string) Trigger {
	return scheduleTrigger{
		cronScheduler: cronScheduler,
		crontab:       crontab,
	}
}

// var _ Trigger = (*webhookTrigger)(nil)

// type webhookTrigger struct {
// 	webhookID string
// }

// func (t webhookTrigger) setup(collector *Collector) error {
// 	normalizedID := url.PathEscape(t.webhookID)

// 	if normalizedID == "" {
// 		return fmt.Errorf("missing webhook ID")
// 	}

// 	// TODO: register webhook handler

// // TODO: add cleanup handler that removes the webhook handler

// 	return nil
// }

// func OnWebhook(webhookID string) Trigger {
// 	return webhookTrigger{
// 		webhookID: webhookID,
// 	}
// }

var _ Trigger = (*changeTrigger)(nil)

type changeTrigger struct {
	eventBroker *events.EventBroker
	pathPattern string
}

func (t changeTrigger) setup(collector *Collector) error {
	if t.pathPattern == "" {
		return fmt.Errorf("missing path pattern")
	}

	go func() {
		for eventBundle := range t.eventBroker.Subscribe(collector.Done()) {
			time := time.Now().UTC().Format(time.RFC3339Nano)

			triggers := make([]any, 0, len(eventBundle))
			for _, event := range eventBundle {
				switch event.Event() {
				case "change":
					path := event.Data()

					if events.MatchEvent(t.pathPattern, path) {
						triggers = append(triggers, map[string]any{
							"type": "change",
							"path": path,
							"time": time,
						})
					}
				}
			}

			collector.scheduler.Schedule(triggers...)
		}
	}()

	return nil
}

func OnChange(eventBroker *events.EventBroker, pathPattern string) Trigger {
	return changeTrigger{
		eventBroker: eventBroker,
		pathPattern: pathPattern,
	}
}
