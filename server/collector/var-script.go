package collector

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

	"github.com/2manyvcos/paranal/server/staticconfig"
	"github.com/go-viper/mapstructure/v2"
	"github.com/google/shlex"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type ScriptConfig struct {
	Command       string             `mapstructure:"command"`
	Env           map[string]*string `mapstructure:"env"`
	Timeout       int                `mapstructure:"timeout"`
	DetailedError bool               `mapstructure:"detailedError"`
}

func VarScript(cfg *staticconfig.Config) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		var err error

		argCount := len(args)
		if argCount < 1 {
			return nil, fmt.Errorf("not enough arguments")
		}
		if argCount > 1 {
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
			stdin = strings.NewReader(v)
		default:
			serializedStdin, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("parsing stdin")
			}
			stdin = bytes.NewReader(serializedStdin)
		}

		strippedArg, err := library.StripJSON(args[0])
		if err != nil {
			return nil, err
		}
		t, err := library.Type(strippedArg)
		if err != nil {
			return nil, err
		}
		var config ScriptConfig
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &config})
		if err != nil {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): %v", t, strippedArg, err.Error()))
		}
		err = decoder.Decode(strippedArg)
		if err != nil {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): %v", t, strippedArg, err.Error()))
		}
		if config.Command == "" {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): missing \"command\"", t, strippedArg))
		}

		command, err := shlex.Split(config.Command)
		if err != nil {
			return nil, fmt.Errorf("parsing command: %v", err)
		}
		if len(command) < 1 {
			return nil, fmt.Errorf("invalid command")
		}
		scriptsPath, err := filepath.Abs(cfg.ScriptDir)
		if err != nil {
			return nil, fmt.Errorf("resolving script: %v", err)
		}
		commandPath, err := filepath.Abs(filepath.Join(cfg.ScriptDir, command[0]))
		if err != nil {
			return nil, fmt.Errorf("resolving script: %v", err)
		}
		if rel, err := filepath.Rel(scriptsPath, commandPath); err != nil {
			return nil, fmt.Errorf("resolving script: %v", err)
		} else if strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || rel == ".." || rel == "." {
			return nil, fmt.Errorf("resolving script: invalid script location")
		}
		var ctx context.Context
		var cancel context.CancelFunc
		if config.Timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
		} else {
			ctx, cancel = context.WithCancel(context.Background())
		}
		defer cancel()
		unsub := signal.Subscribe(cancel)
		defer unsub()
		cmd := exec.CommandContext(ctx, commandPath, command[1:]...)
		cmd.Stdin = stdin
		var stdout strings.Builder
		cmd.Stdout = &stdout
		var stderr strings.Builder
		cmd.Stderr = &stderr
		env := make([]string, 0, len(config.Env))
		for key, value := range config.Env {
			if value != nil {
				env = append(env, key+"="+*value)
			}
		}
		cmd.Env = env
		cmd.Dir = filepath.Dir(commandPath)
		cmd.WaitDelay = 5 * time.Second
		err = cmd.Run()
		if err != nil {
			if config.DetailedError {
				status := -1
				if exitErr, ok := err.(*exec.ExitError); ok {
					status = exitErr.ExitCode()
				}
				return nil, library.ThrowAny(library.NewRuntimeError(map[string]any{
					"status": float64(status),
					"output": stdout.String(),
					"error":  stderr.String(),
				}))
			}
			return nil, fmt.Errorf("running command: %v\n\n%v", err, stderr.String())
		}
		return []any{stdout.String()}, nil
	})
}
