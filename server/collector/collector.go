package collector

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/2manyvcos/paranal/server/config"
	"github.com/2manyvcos/paranal/server/storage"
	"github.com/2manyvcos/paranal/utils"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
)

type Collector struct {
	config    *config.Config
	store     *storage.Storage
	id        string
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
	facts, programErr := c.program.Run([]any{nil}, &jpl.JPLProgramConfig{
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

	states := make([]storage.State, 0, len(facts))

	for _, fact := range facts {
		var f genericFact
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnset: true, Result: &f})
		if err != nil {
			return fmt.Errorf("invalid fact: %w", err)
		}
		err = decoder.Decode(fact)
		if err != nil {
			return fmt.Errorf("invalid fact: %w", err)
		}

		switch f.Kind {
		case "state":
			var s state
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &s})
			if err != nil {
				return fmt.Errorf("invalid state: %w", err)
			}
			err = decoder.Decode(fact)
			if err != nil {
				return fmt.Errorf("invalid state: %w", err)
			}
			err = s.Validate()
			if err != nil {
				return fmt.Errorf("invalid state: %w", err)
			}
			states = append(states, storage.State{
				Registry: s.Registry,
				ID:       s.ID,
				Data:     s.Data,
			})

		default:
			return fmt.Errorf("invalid fact: unsupported kind %q", f.Kind)
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

func New(cfg *config.Config, store *storage.Storage, id string, program jpl.JPLProgram, triggers ...Trigger) (*Collector, error) {
	if id == "" {
		return nil, fmt.Errorf("missing ID")
	}

	var c Collector
	c.config = cfg
	c.store = store
	c.id = id
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

func NewContext(ctx context.Context, cfg *config.Config, store *storage.Storage, id string, program jpl.JPLProgram, triggers ...Trigger) (*Collector, error) {
	c, err := New(cfg, store, id, program, triggers...)
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
