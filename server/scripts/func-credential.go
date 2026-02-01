package scripts

import (
	"errors"
	"fmt"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

func FuncCredential(app *application.App) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		var err error
		var arg any

		if len(args) > 0 {
			arg, err = library.UnwrapValue(args[0])
			if err != nil {
				return nil, err
			}
		} else {
			arg = nil
		}
		t, err := library.Type(arg)
		if err != nil {
			return nil, err
		}
		if t != jpl.JPLT_STRING {
			return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be used as name", string(t), arg))
		}
		name := arg.(string)

		credential, err := app.GetUserCredential(name)
		if errors.Is(err, data.ErrNotFound) {
			return nil, fmt.Errorf("user credential \"%s\" not found", name)
		}
		if err != nil {
			return nil, fmt.Errorf("error fetching user credential \"%s\"", name)
		}
		if credential.Value == "" {
			return nil, fmt.Errorf("user credential \"%s\" has no value", name)
		}
		credential.Value, err = crypto.Decrypt(app.Config.SecretKey, credential.Value)
		if err != nil {
			return nil, fmt.Errorf("error decrypting user credential \"%s\"", name)
		}
		return []any{credential.Value}, nil
	})
}
