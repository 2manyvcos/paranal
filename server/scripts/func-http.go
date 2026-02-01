package scripts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

var jsonRegex = regexp.MustCompile("^application/[^+]*[+]?(json);?.*$")

type Options struct {
	Method  string             `mapstructure:"method"`
	Auth    string             `mapstructure:"auth"`
	Headers map[string]*string `mapstructure:"headers"`
}

func FuncHTTP(app *application.App) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		var err error

		if len(args) < 1 {
			return nil, fmt.Errorf("too view arguments")
		}
		if len(args) > 2 {
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
				return nil, fmt.Errorf("error parsing data")
			}
			body = bytes.NewReader(serializedBody)
			detectedContentType = "application/json"
		}

		unwrappedURL, err := library.UnwrapValue(args[0])
		if err != nil {
			return nil, err
		}
		url, ok := unwrappedURL.(string)
		if !ok {
			return nil, fmt.Errorf("invalid URL")
		}

		var options Options
		options.Method = "GET"
		if body != nil {
			options.Method = "POST"
			options.Headers = map[string]*string{"Content-Type": &detectedContentType}
		}
		if len(args) > 1 {
			var strippedOptions any
			strippedOptions, err = library.StripJSON(args[1])
			if err != nil {
				return nil, err
			}
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &options})
			if err != nil {
				return nil, fmt.Errorf("invalid argument \"options\": %s", err)
			}
			err = decoder.Decode(strippedOptions)
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

		req, err := http.NewRequest(options.Method, url, body)
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
		case data.HTTP_CREDENTIAL_TYPE_BASIC:
			req.SetBasicAuth(credential.Key, credential.Value)
		case data.HTTP_CREDENTIAL_TYPE_BEARER:
			req.Header.Add("Authorization", "Bearer "+credential.Value)
		case data.HTTP_CREDENTIAL_TYPE_HEADER:
			req.Header.Add(credential.Key, credential.Value)
		case data.HTTP_CREDENTIAL_TYPE_QUERY:
			query := req.URL.Query()
			query.Add(credential.Key, credential.Value)
			req.URL.RawQuery = query.Encode()
		default:
			return nil, fmt.Errorf("invalid HTTP credential \"%s\"", credential.Name)
		}
		ctx, cancel := context.WithCancel(req.Context())
		unsub := signal.Subscribe(cancel)
		defer unsub()

		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("error sending HTTP request: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("invalid status %v", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		switch {
		case jsonRegex.MatchString(contentType):
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
