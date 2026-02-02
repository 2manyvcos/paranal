package scripts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/2manyvcos/paranal/server/application"
	"github.com/go-viper/mapstructure/v2"
	"github.com/google/shlex"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type ScriptOptions struct {
	Env     map[string]*string `mapstructure:"env"`
	Timeout int                `mapstructure:"timeout"`
}

func FuncScript(app *application.App) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		var err error

		argCount := len(args)
		if argCount < 1 {
			return nil, fmt.Errorf("not enough arguments")
		} else if argCount > 2 {
			return nil, fmt.Errorf("too many arguments")
		}

		strippedInput, err := library.StripJSON(input)
		if err != nil {
			return nil, err
		}
		var stdin io.Reader
		switch v := strippedInput.(type) {
		case nil:
		case string:
			stdin = bytes.NewReader([]byte(v))
		default:
			serializedStdin, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("error parsing stdin")
			}
			stdin = bytes.NewReader(serializedStdin)
		}

		unwrappedArg, err := library.UnwrapValue(args[0])
		if err != nil {
			return nil, err
		}
		commandStr, ok := unwrappedArg.(string)
		if !ok {
			return nil, fmt.Errorf("invalid command")
		}
		command, err := shlex.Split(commandStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing command: %s", err)
		}
		if len(command) < 1 {
			return nil, fmt.Errorf("invalid command")
		}

		var options ScriptOptions
		if argCount > 1 {
			var err error
			strippedArg, err := library.StripJSON(args[1])
			if err != nil {
				return nil, err
			}
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &options})
			if err != nil {
				return nil, fmt.Errorf("invalid argument \"options\": %s", err)
			}
			err = decoder.Decode(strippedArg)
			if err != nil {
				return nil, fmt.Errorf("invalid argument \"options\": %s", err)
			}
		}

		scriptsPath, err := filepath.Abs(app.Config.ScriptsPath)
		if err != nil {
			return nil, fmt.Errorf("error resolving script: %s", err)
		}
		commandPath, err := filepath.Abs(filepath.Join(app.Config.ScriptsPath, command[0]))
		if err != nil {
			return nil, fmt.Errorf("error resolving script: %s", err)
		}
		if rel, err := filepath.Rel(scriptsPath, commandPath); err != nil {
			return nil, fmt.Errorf("error resolving script: %s", err)
		} else if strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || rel == ".." || rel == "." {
			return nil, fmt.Errorf("error resolving script: invalid script location")
		}
		var ctx context.Context
		var cancel context.CancelFunc
		if options.Timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(options.Timeout)*time.Second)
		} else {
			ctx, cancel = context.WithCancel(context.Background())
		}
		go func() {
			select {
			case <-ctx.Done():
				fmt.Println("context done")
			case <-time.After(5 * time.Second):
				cancel()
			}
		}()
		defer cancel()
		unsub := signal.Subscribe(cancel)
		defer unsub()
		cmd := exec.CommandContext(ctx, commandPath, command[1:]...)
		cmd.Stdin = stdin
		var stdout strings.Builder
		cmd.Stdout = &stdout
		var stderr strings.Builder
		cmd.Stderr = &stderr
		env := make([]string, 0, len(options.Env))
		for key, value := range options.Env {
			if value != nil {
				env = append(env, key+"="+*value)
			}
		}
		cmd.Env = env
		cmd.Dir = filepath.Dir(commandPath)
		cmd.WaitDelay = 5 * time.Second
		err = cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("error running command: %s\n\n%s", err, stderr.String())
		}
		return []any{stdout.String()}, nil
	})
}
