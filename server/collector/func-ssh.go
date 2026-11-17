package collector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Address    string   `mapstructure:"address"`
	HostKeys   []string `mapstructure:"hostKeys"`
	User       string   `mapstructure:"user"`
	Password   string   `mapstructure:"password"`
	PrivateKey string   `mapstructure:"privateKey"`
}

type SSHOptions struct {
	Env           map[string]*string `mapstructure:"env"`
	Timeout       int                `mapstructure:"timeout"`
	DetailedError bool               `mapstructure:"detailedError"`
}

func FuncSSH() jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		var err error

		argCount := len(args)
		if argCount < 2 {
			return nil, fmt.Errorf("not enough arguments")
		} else if argCount > 3 {
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
				return nil, fmt.Errorf("parsing stdin")
			}
			stdin = bytes.NewReader(serializedStdin)
		}

		strippedArg, err := library.StripJSON(args[0])
		if err != nil {
			return nil, err
		}
		var config SSHConfig
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &config})
		if err != nil {
			return nil, fmt.Errorf("invalid argument \"config\": %v", err)
		}
		err = decoder.Decode(strippedArg)
		if err != nil {
			return nil, fmt.Errorf("invalid argument \"config\": %v", err)
		}

		if config.Address == "" {
			return nil, fmt.Errorf("invalid address")
		}
		if len(config.HostKeys) < 1 {
			return nil, fmt.Errorf("missing host key")
		}

		unwrappedArg, err := library.UnwrapValue(args[1])
		if err != nil {
			return nil, err
		}
		command, ok := unwrappedArg.(string)
		if !ok || command == "" {
			return nil, fmt.Errorf("invalid command")
		}

		var options SSHOptions
		if argCount > 2 {
			strippedArg, err = library.StripJSON(args[2])
			if err != nil {
				return nil, err
			}
			decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &options})
			if err != nil {
				return nil, fmt.Errorf("invalid argument \"options\": %v", err)
			}
			err = decoder.Decode(strippedArg)
			if err != nil {
				return nil, fmt.Errorf("invalid argument \"options\": %v", err)
			}
		}

		hostKeys := make([]ssh.PublicKey, len(config.HostKeys))
		for i, key := range config.HostKeys {
			hostKeys[i], _, _, _, err = ssh.ParseAuthorizedKey([]byte(key))
			if err != nil {
				return nil, fmt.Errorf("parsing host key: %v", err)
			}
		}
		sshConfig := &ssh.ClientConfig{
			User: config.User,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				for _, k := range hostKeys {
					if bytes.Equal(key.Marshal(), k.Marshal()) {
						return nil
					}
				}
				return fmt.Errorf("host key mismatch")
			},
		}
		if config.Password != "" {
			sshConfig.Auth = append(sshConfig.Auth, ssh.Password(config.Password))
		}
		if config.PrivateKey != "" {
			signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
			if err != nil {
				return nil, fmt.Errorf("parsing private key: %v", err)
			}
			sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
		}
		client, err := ssh.Dial("tcp", config.Address, sshConfig)
		if err != nil {
			return nil, fmt.Errorf("connecting to host: %v", err)
		}
		defer client.Close()
		session, err := client.NewSession()
		if err != nil {
			return nil, fmt.Errorf("creating remote session: %v", err)
		}
		defer session.Close()
		var ctx context.Context
		var cancel context.CancelFunc
		if options.Timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(options.Timeout)*time.Second)
		} else {
			ctx, cancel = context.WithCancel(context.Background())
		}
		defer cancel()
		unsub := signal.Subscribe(cancel)
		defer unsub()
		go func() {
			<-ctx.Done()
			session.Signal(ssh.SIGINT)
			session.Close()
		}()
		session.Stdin = stdin
		var stdout strings.Builder
		session.Stdout = &stdout
		var stderr strings.Builder
		session.Stderr = &stderr
		for key, value := range options.Env {
			if value == nil {
				err = session.Setenv(key, "")
			} else {
				err = session.Setenv(key, *value)
			}
			if err != nil {
				return nil, fmt.Errorf("setting remote env: %v", err)
			}
		}
		err = session.Run(command)
		if err != nil {
			if options.DetailedError {
				status := -1
				if exitErr, ok := err.(*ssh.ExitError); ok {
					status = exitErr.ExitStatus()
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
