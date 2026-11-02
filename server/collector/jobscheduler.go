package collector

import "sync"

type JobScheduler[TriggerType any] struct {
	triggers                chan TriggerType
	ready, shutdown, done   chan struct{}
	readyOnce, shutdownOnce sync.Once
	job                     func(triggers []TriggerType)
}

// Start tells the scheduler to start running jobs.
func (s *JobScheduler[TriggerType]) Start() {
	s.readyOnce.Do(func() {
		close(s.ready)
	})
}

// Schedule registers a trigger for a job execution.
// Jobs can already be scheduled before calling Start, but only run after calling Start.
func (s *JobScheduler[TriggerType]) Schedule(trigger TriggerType) {
	select {
	case s.triggers <- trigger:
	case <-s.shutdown:
	}
}

// Shutdown tells the scheduler to shutdown.
// If Start has been called, already pending jobs and triggers are processed before the scheduler exits.
// If Start was never called, pending triggers are discarded.
func (s *JobScheduler[TriggerType]) Shutdown() {
	s.shutdownOnce.Do(func() {
		close(s.shutdown)
	})
}

// ShutdownSync calls Shutdown and blocks until the scheduler has exited.
// It must not be called from within the job function, as this would cause a deadlock.
func (s *JobScheduler[TriggerType]) ShutdownSync() {
	s.Shutdown()
	<-s.done
}

func (s *JobScheduler[TriggerType]) handle() {
	triggers := s.triggers
	ready := s.ready
	shutdown := s.shutdown

	var done <-chan struct{}
	var pending []TriggerType

	run := func() {
		batch := pending
		pending = nil

		finished := make(chan struct{})
		done = finished

		go func() {
			defer close(finished)

			s.job(batch)
		}()
	}

	for shutdown != nil || done != nil {
		select {
		case trigger := <-triggers:
			pending = append(pending, trigger)

			if ready == nil && done == nil {
				run()
			}

		case <-ready:
			ready = nil

			if len(pending) > 0 && done == nil {
				run()
			}

		case <-shutdown:
			shutdown = nil
			triggers = nil

			if ready != nil {
				select {
				case <-ready:
					ready = nil

					if len(pending) > 0 && done == nil {
						run()
					}

				default:
				}
			}

		case <-done:
			done = nil

			if len(pending) > 0 {
				run()
			}
		}
	}

	close(s.done)
}

// NewJobScheduler creates a new job scheduler that runs the given job function.
// At most one job runs at a time.
// Triggers that arrive while a job is running are batched and passed together to the next job, which runs immediately after the previous one.
//
// The job function runs in its own goroutine.
func NewJobScheduler[TriggerType any](job func(triggers []TriggerType)) *JobScheduler[TriggerType] {
	var s JobScheduler[TriggerType]
	s.triggers = make(chan TriggerType)
	s.ready = make(chan struct{})
	s.shutdown = make(chan struct{})
	s.done = make(chan struct{})
	s.job = job

	go s.handle()

	return &s
}
