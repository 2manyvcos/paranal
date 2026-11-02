package collector

import "github.com/jplorg/jpl/go/v2/jpl"

func FuncTriggers(triggers []any) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		return []any{triggers}, nil
	})
}
