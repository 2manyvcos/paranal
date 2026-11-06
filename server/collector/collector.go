package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/2manyvcos/paranal/server/config"
	"github.com/2manyvcos/paranal/utils"
	"github.com/jplorg/jpl/go/v2/jpl"
)

type Collector struct {
	config    *config.Config
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

func (p *Collector) run(triggers []any) {
	results, programErr := p.program.Run([]any{nil}, &jpl.JPLProgramConfig{
		Runtime: jpl.JPLRuntimeOptions{
			Vars: map[string]any{
				"http":     FuncHTTP(),
				"script":   FuncScript(p.config),
				"ssh":      FuncSSH(),
				"triggers": FuncTriggers(triggers),
			},
		},
	})
	if programErr != nil {
		// TODO: handle error
		fmt.Println(programErr)
		return
	}

	// TODO: handle results
	json, err := json.Marshal(results)
	fmt.Println(string(json), err)
}

func New(cfg *config.Config, program jpl.JPLProgram, triggers ...Trigger) (*Collector, error) {
	var c Collector
	c.config = cfg
	c.program = program
	c.scheduler = utils.NewJobScheduler(c.run)
	c.done = make(chan struct{})

	ok := true

	if len(triggers) == 0 {
		ok = false
		// TODO: handle error
	}

	for _, trigger := range triggers {
		err := trigger.setup(&c)
		if err != nil {
			ok = false
			// TODO: handle error
			fmt.Println(err)
			break
		}
	}

	if !ok {
		c.Shutdown()

		return nil, fmt.Errorf("setting up collector")
	}

	return &c, nil
}

func NewContext(ctx context.Context, cfg *config.Config, program jpl.JPLProgram, triggers ...Trigger) (*Collector, error) {
	c, err := New(cfg, program, triggers...)
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
