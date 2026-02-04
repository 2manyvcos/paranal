package scripts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/2manyvcos/paranal/utils"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type HTTPOptions struct {
	Method  string             `mapstructure:"method"`
	Auth    string             `mapstructure:"auth"`
	Headers map[string]*string `mapstructure:"headers"`
	Timeout int                `mapstructure:"timeout"`
}

func FuncHTTP(app *application.App) jpl.JPLFunc {
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
		var body io.Reader
		var detectedContentType string
		switch v := strippedInput.(type) {
		case nil:
		case string:
			body = bytes.NewReader([]byte(v))
			detectedContentType = "text/plain"
		default:
			serializedBody, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("error parsing request body")
			}
			body = bytes.NewReader(serializedBody)
			detectedContentType = "application/json"
		}

		unwrappedArg, err := library.UnwrapValue(args[0])
		if err != nil {
			return nil, err
		}
		url, ok := unwrappedArg.(string)
		if !ok {
			return nil, fmt.Errorf("invalid URL")
		}

		var options HTTPOptions
		options.Method = "GET"
		if body != nil {
			options.Method = "POST"
			options.Headers = map[string]*string{"Content-Type": &detectedContentType}
		}
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

		var credential data.HTTPCredential
		if options.Auth != "" {
			credential, err = app.GetHTTPCredential(options.Auth)
			if errors.Is(err, data.ErrNotFound) {
				return nil, fmt.Errorf("HTTP credential \"%s\" not found", options.Auth)
			}
			if err != nil {
				return nil, fmt.Errorf("error fetching HTTP credential \"%s\"", options.Auth)
			}
			if credential.Value == "" {
				return nil, fmt.Errorf("HTTP credential \"%s\" has no value", options.Auth)
			}
			credential.Value, err = crypto.Decrypt(app.Config.SecretKey, credential.Value)
			if err != nil {
				return nil, fmt.Errorf("error decrypting HTTP credential \"%s\"", options.Auth)
			}
		}

		var ctx context.Context
		var cancel context.CancelFunc
		if options.Timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(options.Timeout)*time.Second)
		} else {
			ctx, cancel = context.WithCancel(context.Background())
		}
		unsub := signal.Subscribe(cancel)
		defer unsub()
		req, err := http.NewRequestWithContext(ctx, options.Method, url, body)
		if err != nil {
			return nil, fmt.Errorf("error sending HTTP request: %s", err)
		}
		for key, value := range options.Headers {
			if value == nil {
				req.Header.Del(key)
			} else {
				req.Header.Set(key, *value)
			}
		}
		switch credential.Type {
		case 0:
		case data.HttpCredentialTypeBasic:
			if credential.Key == "" {
				return nil, fmt.Errorf("HTTP credential \"%s\" has no key", credential.Name)
			}
			req.SetBasicAuth(credential.Key, credential.Value)
		case data.HttpCredentialTypeBearer:
			req.Header.Add("Authorization", "Bearer "+credential.Value)
		case data.HttpCredentialTypeHeader:
			if credential.Key == "" {
				return nil, fmt.Errorf("HTTP credential \"%s\" has no key", credential.Name)
			}
			req.Header.Add(credential.Key, credential.Value)
		case data.HttpCredentialTypeQuery:
			if credential.Key == "" {
				return nil, fmt.Errorf("HTTP credential \"%s\" has no key", credential.Name)
			}
			query := req.URL.Query()
			query.Add(credential.Key, credential.Value)
			req.URL.RawQuery = query.Encode()
		default:
			return nil, fmt.Errorf("invalid HTTP credential \"%s\"", credential.Name)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending HTTP request: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("invalid status %v", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		switch {
		case utils.JsonRegex.MatchString(contentType):
			var parsedResponsePayload any
			err = json.NewDecoder(resp.Body).Decode(&parsedResponsePayload)
			if err != nil {
				return nil, fmt.Errorf("error reading HTTP response: %s", err)
			}
			return []any{parsedResponsePayload}, nil

		default:
			responsePayload, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("error reading HTTP response: %s", err)
			}
			return []any{string(responsePayload)}, nil
		}
	})
}
