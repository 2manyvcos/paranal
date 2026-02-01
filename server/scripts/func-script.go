package scripts

import (
	"github.com/2manyvcos/paranal/server/application"
	"github.com/jplorg/jpl/go/v2/jpl"
)

func FuncScript(app *application.App) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		panic("TODO:")
	})
}
