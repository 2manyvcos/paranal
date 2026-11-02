package collector

import (
	"encoding/json"
	"fmt"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/jplorg/jpl/go/v2/jpl"
)

type Collector struct {
	app       *application.App
	program   jpl.JPLProgram
	scheduler *JobScheduler[any]
	done      chan struct{}
}

func (c *Collector) Shutdown() {
	close(c.done)

	c.scheduler.ShutdownSync()
}

func (p *Collector) run(triggers []any) {
	results, programErr := p.program.Run([]any{nil}, &jpl.JPLProgramConfig{
		Runtime: jpl.JPLRuntimeOptions{
			Vars: map[string]any{
				"http":     FuncHTTP(),
				"script":   FuncScript(p.app),
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

func NewCollector(app *application.App, program jpl.JPLProgram, triggers ...Trigger) (*Collector, error) {
	var c Collector
	c.app = app
	c.program = program
	c.scheduler = NewJobScheduler(c.run)
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

		return nil, fmt.Errorf("setting up collector failed")
	}

	c.scheduler.Start()

	return &c, nil
}
