package scripts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
	"golang.org/x/crypto/ssh"
)

type SSHHost struct {
	Address  string   `mapstructure:"address"`
	HostKeys []string `mapstructure:"hostKeys"`
	Auth     string   `mapstructure:"auth"`
}

type SSHOptions struct {
	Env     map[string]*string `mapstructure:"env"`
	Timeout int                `mapstructure:"timeout"`
}

func FuncSSH(app *application.App) jpl.JPLFunc {
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
				return nil, fmt.Errorf("error parsing stdin")
			}
			stdin = bytes.NewReader(serializedStdin)
		}

		strippedArg, err := library.StripJSON(args[0])
		if err != nil {
			return nil, err
		}
		var host SSHHost
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &host})
		if err != nil {
			return nil, fmt.Errorf("invalid argument \"host\": %s", err)
		}
		err = decoder.Decode(strippedArg)
		if err != nil {
			return nil, fmt.Errorf("invalid argument \"host\": %s", err)
		}

		if host.Address == "" {
			return nil, fmt.Errorf("invalid address")
		}
		if len(host.HostKeys) < 1 {
			return nil, fmt.Errorf("host has no host key")
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
				return nil, fmt.Errorf("invalid argument \"options\": %s", err)
			}
			err = decoder.Decode(strippedArg)
			if err != nil {
				return nil, fmt.Errorf("invalid argument \"options\": %s", err)
			}
		}

		var credential data.SSHCredential
		if host.Auth != "" {
			credential, err = app.GetSSHCredential(host.Auth)
			if errors.Is(err, data.ErrNotFound) {
				return nil, fmt.Errorf("SSH credential \"%s\" not found", host.Auth)
			}
			if err != nil {
				return nil, fmt.Errorf("error fetching SSH credential \"%s\"", host.Auth)
			}
			if credential.User == "" {
				return nil, fmt.Errorf("SSH credential \"%s\" has no user", host.Auth)
			}
			if credential.Password != "" {
				credential.Password, err = crypto.Decrypt(app.Config.SecretKey, credential.Password)
				if err != nil {
					return nil, fmt.Errorf("error decrypting SSH credential \"%s\"", host.Auth)
				}
			}
			if credential.PrivateKey != "" {
				credential.PrivateKey, err = crypto.Decrypt(app.Config.SecretKey, credential.PrivateKey)
				if err != nil {
					return nil, fmt.Errorf("error decrypting SSH credential \"%s\"", host.Auth)
				}
			}
		}

		hostKeys := make([]ssh.PublicKey, len(host.HostKeys))
		for i, key := range host.HostKeys {
			hostKeys[i], _, _, _, err = ssh.ParseAuthorizedKey([]byte(key))
			if err != nil {
				return nil, fmt.Errorf("error parsing public key: %s", err)
			}
		}
		sshConfig := &ssh.ClientConfig{
			User: credential.User,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				for _, k := range hostKeys {
					if bytes.Equal(key.Marshal(), k.Marshal()) {
						return nil
					}
				}
				return fmt.Errorf("ssh: host key mismatch")
			},
		}
		if credential.Password != "" {
			sshConfig.Auth = append(sshConfig.Auth, ssh.Password(credential.Password))
		}
		if credential.PrivateKey != "" {
			signer, err := ssh.ParsePrivateKey([]byte(credential.PrivateKey))
			if err != nil {
				return nil, fmt.Errorf("error parsing private key: %s", err)
			}
			sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
		}
		client, err := ssh.Dial("tcp", host.Address, sshConfig)
		if err != nil {
			return nil, fmt.Errorf("error connecting to host: %s", err)
		}
		defer client.Close()
		session, err := client.NewSession()
		if err != nil {
			return nil, fmt.Errorf("error creating remote session: %s", err)
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
				return nil, fmt.Errorf("error setting remote env: %s", err)
			}
		}
		err = session.Run(command)
		if err != nil {
			return nil, fmt.Errorf("error running command: %s\n\n%s", err, stderr.String())
		}
		return []any{stdout.String()}, nil
	})
}
