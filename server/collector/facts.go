package collector

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
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

func parseState(fact any) (state, error) {
	var s state
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &s})
	if err != nil {
		return state{}, err
	}
	err = decoder.Decode(fact)
	if err != nil {
		return state{}, err
	}

	if s.Registry == "" {
		return state{}, fmt.Errorf("missing \"registry\"")
	}

	if s.ID == "" {
		return state{}, fmt.Errorf("missing \"id\"")
	}

	if s.Data == nil {
		return state{}, fmt.Errorf("missing \"data\"")
	}
	if _, ok := s.Data.(map[string]any); !ok {
		return state{}, fmt.Errorf("\"data\" must be an object")
	}

	return s, nil
}
