package scripts

import (
	"github.com/2manyvcos/paranal/server/application"
	gojpl "github.com/jplorg/jpl/go/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
)

func Config(app *application.App, serviceID string) *jpl.JPLInterpreterConfig {
	return &jpl.JPLInterpreterConfig{
		Runtime: jpl.JPLRuntimeOptions{
			Vars: map[string]any{
				"script":     FuncScript(app),
				"ssh":        FuncSSH(app),
				"http":       FuncHTTP(app),
				"credential": FuncCredential(app),
				"service":    FuncService(app, serviceID),
			},
		},
	}
}

func Parse(app *application.App, serviceID string, script string) (jpl.JPLProgram, error) {
	return gojpl.Parse(script, Config(app, serviceID))
}

func Run(app *application.App, serviceID string, script string) ([]any, error) {
	return gojpl.Run(script, []any{nil}, Config(app, serviceID))
}
