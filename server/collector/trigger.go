package collector

import (
	"time"

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

// var _ Trigger = (*eventTrigger)(nil)

// type eventTrigger struct {
// 	eventPath string
// }

// func (t eventTrigger) setup(collector *Collector) error {
// if normalizedEventPath == "" {
// 		return fmt.Errorf("missing event path")
// 	}

// 	// TODO: register event handler (also handle wildcards)

// // TODO: add cleanup handler that removes the event handler

// 	return nil
// }

// func OnEvent(eventPath string) Trigger {
// 	return eventTrigger{
// 		eventPath: eventPath,
// 	}
// }
