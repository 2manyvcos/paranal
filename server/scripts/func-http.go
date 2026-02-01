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
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

var jsonRegex = regexp.MustCompile("^application/[^+]*[+]?(json);?.*$")

func FuncHTTP(app *application.App) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error) {
		var err error
		var arg any

		arg, err = library.StripJSON(input)
		if err != nil {
			return nil, err
		}
		t, err := library.Type(arg)
		if err != nil {
			return nil, err
		}
		var body io.Reader
		var detectedContentType string
		switch t {
		case jpl.JPLT_NULL:
		case jpl.JPLT_STRING:
			body = bytes.NewReader([]byte(arg.(string)))
			detectedContentType = "text/plain"
		default:
			serializedBody, err := json.Marshal(arg)
			if err != nil {
				return nil, fmt.Errorf("error parsing data")
			}
			body = bytes.NewReader(serializedBody)
			detectedContentType = "application/json"
		}

		if len(args) > 0 {
			arg, err = library.UnwrapValue(args[0])
			if err != nil {
				return nil, err
			}
		} else {
			arg = nil
		}
		t, err = library.Type(arg)
		if err != nil {
			return nil, err
		}
		if t != jpl.JPLT_STRING {
			return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be used as URL", string(t), arg))
		}
		url := arg.(string)
		if url == "" {
			return nil, fmt.Errorf("invalid URL")
		}

		if len(args) > 1 {
			arg, err = library.StripJSON(args[1])
			if err != nil {
				return nil, err
			}
		} else {
			arg = nil
		}
		t, err = library.Type(arg)
		if err != nil {
			return nil, err
		}
		var options struct {
			Method  string
			Auth    string
			Headers map[string]string
		}
		options.Method = "GET"
		options.Headers = make(map[string]string)
		if body != nil {
			options.Method = "POST"
			options.Headers = map[string]string{"Content-Type": detectedContentType}
		}
		switch t {
		case jpl.JPLT_NULL:
		case jpl.JPLT_OBJECT:
			optionMap := arg.(map[string]any)

			option := optionMap["method"]
			to, err := library.Type(option)
			if err != nil {
				return nil, err
			}
			switch to {
			case jpl.JPLT_NULL:
			case jpl.JPLT_STRING:
				options.Method = option.(string)
			default:
				return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be used as options.method", string(to), option))
			}

			option = optionMap["auth"]
			to, err = library.Type(option)
			if err != nil {
				return nil, err
			}
			switch to {
			case jpl.JPLT_NULL:
			case jpl.JPLT_STRING:
				options.Auth = option.(string)
			default:
				return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be used as options.auth", string(to), option))
			}

			option = optionMap["headers"]
			to, err = library.Type(option)
			if err != nil {
				return nil, err
			}
			switch to {
			case jpl.JPLT_NULL:
			case jpl.JPLT_OBJECT:
				for key, header := range option.(map[string]any) {
					th, err := library.Type(header)
					if err != nil {
						return nil, err
					}
					switch th {
					case jpl.JPLT_NULL:
						delete(options.Headers, key)
					case jpl.JPLT_STRING:
						options.Headers[key] = header.(string)
					default:
						return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be used as header", string(to), option))
					}
				}
			default:
				return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be used as options.headers", string(to), option))
			}
		default:
			return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be used as options", string(t), arg))
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
			req.Header.Set(key, value)
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
