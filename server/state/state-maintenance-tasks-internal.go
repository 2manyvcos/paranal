package state

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/2manyvcos/paranal/server/schema"
	"github.com/go-co-op/gocron/v2"
)

func (s *State) setupMaintenanceTasks() error {
	s.maintenanceTasks = make(map[string]maintenanceTask)

	if err := s.setupManualMaintenanceTask("Run all service scripts", s.RunServiceScripts); err != nil {
		return err
	}

	if err := s.setupScheduledMaintenanceTask("Cleanup database", gocron.CronJob("0 0 * * *", false), s.app.CleanupDatabase); err != nil {
		return err
	}

	return nil
}

func (s *State) setupScheduledMaintenanceTask(name string, schedule gocron.JobDefinition, fn func()) error {
	var task scheduledMaintenanceTask
	if job, err := s.app.Scheduler.NewJob(
		schedule,
		gocron.NewTask(newMaintenanceTaskRunner, s, name, &task, fn),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	); err == nil {
		task.job = job
	} else {
		log.Printf("Error scheduling maintenance task \"%s\" - %s\n", name, err)
		return err
	}
	s.maintenanceTasks[name] = &task
	return nil
}

type scheduledMaintenanceTask struct {
	lock    sync.RWMutex
	job     gocron.Job
	running bool
}

func (t *scheduledMaintenanceTask) State() (lastRun *time.Time, running bool, nextRun *time.Time) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if lastJobRun, err := t.job.LastRun(); err == nil && !lastJobRun.IsZero() {
		lastRun = &lastJobRun
	}
	running = t.running
	if nextJobRun, err := t.job.NextRun(); err == nil && !nextJobRun.IsZero() {
		nextRun = &nextJobRun
	}
	return
}

func (t *scheduledMaintenanceTask) Run() error {
	return t.job.RunNow()
}

func newMaintenanceTaskRunner(s *State, name string, task *scheduledMaintenanceTask, fn func()) {
	task.lock.Lock()
	task.running = true
	s.app.PublishClientEvent(schema.NewUpdateEvent("/v1/maintenance-tasks/" + url.PathEscape(name)))
	task.lock.Unlock()

	log.Printf("Running maintenance task \"%s\"\n", name)

	fn()

	task.lock.Lock()
	task.running = false
	s.app.PublishClientEvent(schema.NewUpdateEvent("/v1/maintenance-tasks/" + url.PathEscape(name)))
	task.lock.Unlock()
}

func (s *State) setupManualMaintenanceTask(name string, fn func()) error {
	var task manualMaintenanceTask
	task.run = func() error {
		go func() {
			task.lock.Lock()
			task.lastRun = time.Now()
			task.running = true
			s.app.PublishClientEvent(schema.NewUpdateEvent("/v1/maintenance-tasks/" + url.PathEscape(name)))
			task.lock.Unlock()

			log.Printf("Running maintenance task \"%s\"\n", name)

			fn()

			task.lock.Lock()
			task.running = false
			s.app.PublishClientEvent(schema.NewUpdateEvent("/v1/maintenance-tasks/" + url.PathEscape(name)))
			task.lock.Unlock()
		}()
		return nil
	}
	s.maintenanceTasks[name] = &task
	return nil
}

type manualMaintenanceTask struct {
	lock    sync.RWMutex
	run     func() error
	running bool
	lastRun time.Time
}

func (t *manualMaintenanceTask) State() (lastRun *time.Time, running bool, nextRun *time.Time) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if lastJobRun := t.lastRun; !lastJobRun.IsZero() {
		lastRun = &lastJobRun
	}
	running = t.running
	return
}

func (t *manualMaintenanceTask) Run() error {
	return t.run()
}
