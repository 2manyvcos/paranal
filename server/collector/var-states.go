package collector

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/server/storage"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type StatesConfig struct {
	Registry string `mapstructure:"registry"`
}

func VarStates(store *storage.Storage) jpl.JPLFunc {
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
		var config StatesConfig
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

		records, err := store.ListStates(config.Registry)
		if err != nil {
			log.Printf("loading records: %v\n", err)
			return nil, errors.New(http.StatusText(http.StatusInternalServerError))
		}

		results := make([]any, len(records))
		for i, record := range records {
			result := make(map[string]any, 3)
			result["registry"] = record.Registry
			result["id"] = record.ID
			result["data"] = record.Data
			results[i] = result
		}

		return []any{results}, nil
	})
}
