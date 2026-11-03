package collector

import (
	"fmt"
	"time"

	"github.com/2manyvcos/paranal/server/event"
	"github.com/go-co-op/gocron/v2"
)

type Trigger interface {
	setup(collector *Collector) error
}

var _ Trigger = (*startTrigger)(nil)

type startTrigger struct{}

func (t startTrigger) setup(collector *Collector) error {
	collector.scheduler.Schedule(map[string]any{
		"type": "start",
		"time": time.Now().UTC().Format(time.RFC3339Nano),
	})

	return nil
}

func OnStart() Trigger {
	return startTrigger{}
}

var _ Trigger = (*scheduleTrigger)(nil)

type scheduleTrigger struct {
	crontab string
}

func (t scheduleTrigger) setup(collector *Collector) error {
	job, err := collector.app.CronScheduler.NewJob(
		gocron.CronJob(t.crontab, true),
		gocron.NewTask(func() {
			collector.scheduler.Schedule(map[string]any{
				"type": "schedule",
				"time": time.Now().UTC().Format(time.RFC3339Nano),
			})
		}),
	)
	if err != nil {
		return err
	}

	go func() {
		<-collector.done
		collector.app.CronScheduler.RemoveJob(job.ID())
	}()

	return nil
}

func OnSchedule(crontab string) Trigger {
	return scheduleTrigger{
		crontab: crontab,
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
	pathPattern string
}

func (t changeTrigger) setup(collector *Collector) error {
	if t.pathPattern == "" {
		return fmt.Errorf("missing path pattern")
	}

	go func() {
		for e := range collector.app.Events.Subscribe(collector.done) {
			switch e.Event() {
			case "change":
				path := e.Data()

				if event.MatchEvent(t.pathPattern, path) {
					collector.scheduler.Schedule(map[string]any{
						"type": "change",
						"path": path,
						"time": time.Now().UTC().Format(time.RFC3339Nano),
					})
				}
			}
		}
	}()

	return nil
}

func OnChange(pathPattern string) Trigger {
	return changeTrigger{
		pathPattern: pathPattern,
	}
}
