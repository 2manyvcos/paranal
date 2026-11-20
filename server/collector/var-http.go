package collector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/2manyvcos/paranal/utils"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type HTTPConfig struct {
	URL           string             `mapstructure:"url"`
	Method        string             `mapstructure:"method"`
	Headers       map[string]*string `mapstructure:"headers"`
	Timeout       int                `mapstructure:"timeout"`
	DetailedError bool               `mapstructure:"detailedError"`
}

func VarHTTP() jpl.JPLFunc {
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
				return nil, fmt.Errorf("parsing request body")
			}
			body = bytes.NewReader(serializedBody)
			detectedContentType = "application/json"
		}

		strippedArg, err := library.StripJSON(args[0])
		if err != nil {
			return nil, err
		}
		t, err := library.Type(strippedArg)
		if err != nil {
			return nil, err
		}
		var config HTTPConfig
		config.Method = "GET"
		if body != nil {
			config.Method = "POST"
			config.Headers = map[string]*string{"Content-Type": &detectedContentType}
		}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: &config})
		if err != nil {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): %v", t, strippedArg, err.Error()))
		}
		err = decoder.Decode(strippedArg)
		if err != nil {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): %v", t, strippedArg, err.Error()))
		}
		if config.URL == "" {
			return nil, library.ThrowAny(library.NewRuntimeError("invalid argument 0 (%s %*<100v): missing \"url\"", t, strippedArg))
		}

		var ctx context.Context
		var cancel context.CancelFunc
		if config.Timeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
		} else {
			ctx, cancel = context.WithCancel(context.Background())
		}
		unsub := signal.Subscribe(cancel)
		defer unsub()
		req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, body)
		if err != nil {
			return nil, fmt.Errorf("sending HTTP request: %v", err)
		}
		for key, value := range config.Headers {
			if value == nil {
				req.Header.Del(key)
			} else {
				req.Header.Set(key, *value)
			}
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("sending HTTP request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			if resp.StatusCode >= 400 && config.DetailedError {
				message, _ := getResponsePayload(resp)
				return nil, library.ThrowAny(library.NewRuntimeError(map[string]any{
					"status": float64(resp.StatusCode),
					"body":   message,
				}))
			}
			return nil, fmt.Errorf("status %v", resp.StatusCode)
		}

		result, err := getResponsePayload(resp)
		if err != nil {
			return nil, fmt.Errorf("reading HTTP response: %v", err)
		}

		return []any{result}, nil
	})
}

func getResponsePayload(resp *http.Response) (any, error) {
	contentType := resp.Header.Get("Content-Type")

	switch {
	case utils.IsJSON(contentType):
		var parsedResponsePayload any
		err := json.NewDecoder(resp.Body).Decode(&parsedResponsePayload)
		if err != nil {
			return nil, err
		}
		return parsedResponsePayload, nil

	default:
		responsePayload, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return string(responsePayload), nil
	}
}
