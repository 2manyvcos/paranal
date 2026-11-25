package collector

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/server/db"
	"github.com/2manyvcos/paranal/server/storage"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type StateConfig struct {
	Registry string `mapstructure:"registry"`
	ID       string `mapstructure:"id"`
}

func VarState(store *storage.Storage) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		var err error

		argCount := len(args)
		if argCount < 1 {
			return nil, fmt.Errorf("not enough arguments")
		}
		if argCount > 1 {
			return nil, fmt.Errorf("too many arguments")
		}

		strippedArg, err := library.StripJSON(args[0])
		if err != nil {
			return nil, err
		}
		t, err := library.Type(strippedArg)
		if err != nil {
			return nil, err
		}
		var config StateConfig
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &config})
		if err != nil {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): %v", t, strippedArg, err.Error()))
		}
		err = decoder.Decode(strippedArg)
		if err != nil {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): %v", t, strippedArg, err.Error()))
		}
		if config.Registry == "" {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): missing \"registry\"", t, strippedArg))
		}
		if config.ID == "" {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): missing \"id\"", t, strippedArg))
		}

		record, err := store.GetState(config.Registry, config.ID)
		if errors.Is(err, db.ErrNotFound) {
			return []any{nil}, nil
		}
		if err != nil {
			log.Printf("loading record: %v\n", err)
			return nil, errors.New(http.StatusText(http.StatusInternalServerError))
		}

		result := make(map[string]any, 3)
		result["registry"] = record.Registry
		result["id"] = record.ID
		result["data"] = record.Data

		return []any{result}, nil
	})
}
