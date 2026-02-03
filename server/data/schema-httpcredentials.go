package data

import (
	"fmt"
)

const (
	HTTP_CREDENTIAL_TYPE_BASIC  = 1
	HTTP_CREDENTIAL_TYPE_BEARER = 2
	HTTP_CREDENTIAL_TYPE_HEADER = 3
	HTTP_CREDENTIAL_TYPE_QUERY  = 4
)

var (
	HTTP_CREDENTIAL_TYPE_NAMES = map[int]string{
		HTTP_CREDENTIAL_TYPE_BASIC:  "basic",
		HTTP_CREDENTIAL_TYPE_BEARER: "bearer",
		HTTP_CREDENTIAL_TYPE_HEADER: "header",
		HTTP_CREDENTIAL_TYPE_QUERY:  "query",
	}
	HTTP_CREDENTIAL_TYPE_CODES map[string]int
)

func init() {
	HTTP_CREDENTIAL_TYPE_CODES = make(map[string]int, len(HTTP_CREDENTIAL_TYPE_NAMES))
	for code, name := range HTTP_CREDENTIAL_TYPE_NAMES {
		HTTP_CREDENTIAL_TYPE_CODES[name] = code
	}
}

type HTTPCredentialTable struct{ TableHeader }

var HTTPCredentials Table[HTTPCredential] = HTTPCredentialTable{
	TableHeader{
		Name:   "httpcredentials",
		Fields: httpCredentialFields,
	},
}

func (t HTTPCredentialTable) NewRecord() TableRecord[HTTPCredential] {
	return new(HTTPCredential)
}

type HTTPCredential struct {
	Name  string
	Type  int
	Key   string
	Value string
}

var httpCredentialFields = []string{"name", "type", "key", "value"}

func (r *HTTPCredential) RecordFields() []any {
	return []any{&r.Name, &r.Type, &r.Key, &r.Value}
}

func (r *HTTPCredential) Dataset() HTTPCredential {
	return *r
}

var _ Upsertable[HTTPCredential] = HTTPCredential{}

func (r HTTPCredential) Table() Table[HTTPCredential] {
	return HTTPCredentials
}

func (r HTTPCredential) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	switch r.Type {
	case HTTP_CREDENTIAL_TYPE_BASIC:
		if r.Key == "" {
			return fmt.Errorf("invalid user")
		}
		if r.Value == "" {
			return fmt.Errorf("invalid password")
		}
	case HTTP_CREDENTIAL_TYPE_BEARER:
		if r.Value == "" {
			return fmt.Errorf("invalid token")
		}
	case HTTP_CREDENTIAL_TYPE_HEADER:
		if r.Key == "" {
			return fmt.Errorf("invalid header name")
		}
		if r.Value == "" {
			return fmt.Errorf("invalid header value")
		}
	case HTTP_CREDENTIAL_TYPE_QUERY:
		if r.Key == "" {
			return fmt.Errorf("invalid query key")
		}
		if r.Value == "" {
			return fmt.Errorf("invalid query value")
		}
	default:
		return fmt.Errorf("invalid type")
	}
	return nil
}

func (r HTTPCredential) InsertableNames() []string {
	return httpCredentialInsertables
}

var httpCredentialInsertables = []string{"name", "type", "key", "value"}

func (r HTTPCredential) Insertables() []any {
	return []any{r.Name, r.Type, r.Key, r.Value}
}

func (r HTTPCredential) UpdatableNames() []string {
	return httpCredentialUpdatables
}

var httpCredentialUpdatables = []string{"name", "type", "key", "value"}

func (r HTTPCredential) Updatables() []any {
	return []any{r.Name, r.Type, r.Key, r.Value}
}

type HTTPCredentialName string

var _ Identifier[HTTPCredential] = HTTPCredentialName("")

func (i HTTPCredentialName) Table() Table[HTTPCredential] {
	return HTTPCredentials
}

func (i HTTPCredentialName) Valid() error {
	if i == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (i HTTPCredentialName) IDNames() []string {
	return HTTPCredentialNameIDs
}

var HTTPCredentialNameIDs = []string{"name"}

func (i HTTPCredentialName) IDs() []any {
	return []any{string(i)}
}
