package collector

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/2manyvcos/paranal/server/staticconfig"
	"github.com/2manyvcos/paranal/server/storage"
	"github.com/2manyvcos/paranal/spec"
	collector_spec "github.com/2manyvcos/paranal/spec/collector"
	"github.com/2manyvcos/paranal/utils"
	"github.com/jplorg/jpl/go/v2/jpl"
)

type Collector struct {
	config    *staticconfig.Config
	store     *storage.Storage
	id        string
	parser    collector_spec.CollectorSpecParser
	program   jpl.JPLProgram
	scheduler *utils.JobScheduler[any]
	done      chan struct{}
	doneOnce  sync.Once
}

func (c *Collector) Start() {
	c.scheduler.Start()
}

func (c *Collector) Done() <-chan struct{} {
	return c.done
}

func (c *Collector) Shutdown() {
	c.doneOnce.Do(func() {
		close(c.done)
	})

	c.scheduler.ShutdownSync()
}

func (c *Collector) run(triggers []any) error {
	results, programErr := c.program.Run([]any{nil}, &jpl.JPLProgramConfig{
		Runtime: jpl.JPLRuntimeOptions{
			Vars: map[string]any{
				"http":     FuncHTTP(),
				"script":   FuncScript(c.config),
				"ssh":      FuncSSH(),
				"triggers": FuncTriggers(triggers),
			},
		},
	})
	if programErr != nil {
		return programErr
	}

	states := make([]collector_spec.State, 0, len(results))
	for _, result := range results {
		fact, err := c.parser.ParseFact(result)
		if err != nil {
			return fmt.Errorf("parsing fact: %w", err)
		}

		switch fact.Kind() {
		case "state":
			states = append(states, fact.(collector_spec.State))

		default:
			return fmt.Errorf("unsupported kind %q", fact.Kind())
		}
	}

	err := c.store.PutCollectorStates(c.id, states)
	if err != nil {
		return err
	}

	return nil
}

func (c *Collector) handle(triggers []any) {
	log.Printf("running collector %q\n", c.id)

	err := c.run(triggers)
	if err != nil {
		log.Printf("running collector %q: %v\n", c.id, err)
	}
}

func New(cfg *staticconfig.Config, store *storage.Storage, version string, id string, program jpl.JPLProgram, triggers ...Trigger) (*Collector, error) {
	if id == "" {
		return nil, fmt.Errorf("missing ID")
	}

	parser, err := spec.CollectionParser(version)
	if err != nil {
		return nil, err
	}

	var c Collector
	c.config = cfg
	c.store = store
	c.id = id
	c.parser = parser
	c.program = program
	c.scheduler = utils.NewJobScheduler(c.handle)
	c.done = make(chan struct{})

	var triggerErrors []error

	if len(triggers) == 0 {
		triggerErrors = append(triggerErrors, fmt.Errorf("no triggers provided"))
	}

	for _, trigger := range triggers {
		err := trigger.setup(&c)
		if err != nil {
			triggerErrors = append(triggerErrors, fmt.Errorf("setting up trigger: %w", err))
			break
		}
	}

	if len(triggerErrors) > 0 {
		c.Shutdown()

		return nil, errors.Join(triggerErrors...)
	}

	return &c, nil
}

func NewContext(ctx context.Context, cfg *staticconfig.Config, store *storage.Storage, version string, id string, program jpl.JPLProgram, triggers ...Trigger) (*Collector, error) {
	c, err := New(cfg, store, version, id, program, triggers...)
	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-ctx.Done():
			c.Shutdown()

		case <-c.Done():
		}
	}()

	return c, nil
}
