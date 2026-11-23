package collector

import (
	"fmt"

	"github.com/2manyvcos/paranal/utils"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

func VarScrapeLogo() jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		var err error

		argCount := len(args)
		if argCount > 0 {
			return nil, fmt.Errorf("too many arguments")
		}

		unwrappedInput, err := library.UnwrapValue(input)
		if err != nil {
			return nil, err
		}
		t, err := library.Type(unwrappedInput)
		if err != nil {
			return nil, err
		}
		url, ok := unwrappedInput.(string)
		if !ok {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid input (%s %*<100v): expected string", t, unwrappedInput))
		}

		logo, err := utils.ScrapeLogo(url)
		if err != nil {
			return nil, err
		}

		if logo == "" {
			return []any{nil}, nil
		}

		return []any{logo}, nil
	})
}
