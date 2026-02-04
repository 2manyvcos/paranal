package data

import (
	"fmt"
)

const (
	HttpCredentialTypeBasic = iota + 1
	HttpCredentialTypeBearer
	HttpCredentialTypeHeader
	HttpCredentialTypeQuery
)

var (
	HttpCredentialTypeNames = map[int]string{
		HttpCredentialTypeBasic:  "basic",
		HttpCredentialTypeBearer: "bearer",
		HttpCredentialTypeHeader: "header",
		HttpCredentialTypeQuery:  "query",
	}
	HttpCredentialTypeCodes map[string]int
)

func init() {
	HttpCredentialTypeCodes = make(map[string]int, len(HttpCredentialTypeNames))
	for code, name := range HttpCredentialTypeNames {
		HttpCredentialTypeCodes[name] = code
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

func (r HTTPCredential) TableType() Table[HTTPCredential] { return nil }

func (r HTTPCredential) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	switch r.Type {
	case HttpCredentialTypeBasic:
		if r.Key == "" {
			return fmt.Errorf("invalid user")
		}
		if r.Value == "" {
			return fmt.Errorf("invalid password")
		}
	case HttpCredentialTypeBearer:
		if r.Value == "" {
			return fmt.Errorf("invalid token")
		}
	case HttpCredentialTypeHeader:
		if r.Key == "" {
			return fmt.Errorf("invalid header name")
		}
		if r.Value == "" {
			return fmt.Errorf("invalid header value")
		}
	case HttpCredentialTypeQuery:
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
