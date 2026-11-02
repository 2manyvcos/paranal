package collector

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
)

type Environment interface {
	CronScheduler() gocron.Scheduler
}

type Collector struct {
	scheduler *JobScheduler[any]
	done      chan struct{}
}

func (c *Collector) Shutdown() {
	close(c.done)

	c.scheduler.ShutdownSync()
}

func (p *Collector) run(triggers []any) {
	// results, err := p.program.Run([]any{nil}, nil)

	// TODO:
}

func NewCollector(env Environment, program jpl.JPLProgram, triggers ...Trigger) (*Collector, error) {
	var c Collector
	c.scheduler = NewJobScheduler(c.run)
	c.done = make(chan struct{})

	for _, trigger := range triggers {
		err := trigger.setup(&c, env)
		if err != nil {
			return nil, err
		}
	}

	c.scheduler.Start()

	return &c, nil
}
