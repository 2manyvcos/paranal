package collector

import (
	"errors"
	"fmt"
)

type genericFact struct {
	Kind string `mapstructure:"kind"`
}

type state struct {
	genericFact `mapstructure:",squash"`
	Registry    string `mapstructure:"registry"`
	ID          string `mapstructure:"id"`
	Data        any    `mapstructure:"data"`
}

func (s state) Validate() error {
	var e []error

	if s.Registry == "" {
		e = append(e, fmt.Errorf("missing \"registry\""))
	}

	if s.ID == "" {
		e = append(e, fmt.Errorf("missing \"id\""))
	}

	if s.Data == nil {
		e = append(e, fmt.Errorf("missing \"data\""))
	} else if _, ok := s.Data.(map[string]any); !ok {
		e = append(e, fmt.Errorf("\"data\" must be an object"))
	}

	return errors.Join(e...)
}
