package collector

import (
	"fmt"

	"github.com/jplorg/jpl/go/v2/jpl"
)

func FuncTriggers(triggers []any) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		argCount := len(args)
		if argCount > 0 {
			return nil, fmt.Errorf("too many arguments")
		}

		return []any{triggers}, nil
	})
}
