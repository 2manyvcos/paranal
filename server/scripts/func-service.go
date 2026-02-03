package scripts

import (
	"errors"
	"fmt"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/jplorg/jpl/go/v2/jpl"
)

func FuncService(app *application.App, serviceID string) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		var err error

		if argCount := len(args); argCount > 0 {
			return nil, fmt.Errorf("too many arguments")
		}

		service, err := app.GetService(serviceID)
		if errors.Is(err, data.ErrNotFound) {
			return nil, fmt.Errorf("service not found")
		}
		if err != nil {
			return nil, fmt.Errorf("error fetching service")
		}
		return []any{map[string]any{
			"id":          service.ID,
			"name":        service.Name,
			"description": service.Description,
			"logo":        service.Logo,
			"url":         service.URL,
		}}, nil
	})
}
