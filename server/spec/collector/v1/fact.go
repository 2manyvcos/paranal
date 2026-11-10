package collector_spec_v1

import (
	"fmt"

	collector_spec "github.com/2manyvcos/paranal/server/spec/collector"
	"github.com/go-viper/mapstructure/v2"
)

type fact struct {
	Kind    string         `mapstructure:"kind"`
	Payload map[string]any `mapstructure:",remain"`
}

func (CollectorParserV1) ParseFact(input any) (collector_spec.Fact, error) {
	var f fact
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &f})
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(input)
	if err != nil {
		return nil, err
	}

	switch f.Kind {
	case "state":
		result, err := parseState(f.Payload)
		if err != nil {
			return nil, fmt.Errorf("parsing state: %w", err)
		}
		return result, nil

	default:
		return nil, fmt.Errorf("unsupported kind %q", f.Kind)
	}
}

type state struct {
	Registry string `mapstructure:"registry"`
	ID       string `mapstructure:"id"`
	Data     any    `mapstructure:"data"`
}

func parseState(input any) (collector_spec.Fact, error) {
	var s state
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &s})
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(input)
	if err != nil {
		return nil, err
	}

	var result collector_spec.State

	if s.Registry == "" {
		return nil, fmt.Errorf("missing \"registry\"")
	}
	result.Registry = s.Registry

	if s.ID == "" {
		return nil, fmt.Errorf("missing \"id\"")
	}
	result.ID = s.ID

	if s.Data == nil {
		return nil, fmt.Errorf("missing \"data\"")
	}
	if _, ok := s.Data.(map[string]any); !ok {
		return nil, fmt.Errorf("\"data\" must be an object")
	}
	result.Data = s.Data

	return result, nil
}
