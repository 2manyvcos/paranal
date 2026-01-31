package data

import (
	"fmt"
)

const HTTP_CREDENTIAL_TYPE_BASIC = 1
const HTTP_CREDENTIAL_TYPE_BEARER = 2
const HTTP_CREDENTIAL_TYPE_HEADER = 3
const HTTP_CREDENTIAL_TYPE_QUERY = 4

var HTTP_CREDENTIAL_TYPE_NAMES = map[int]string{
	HTTP_CREDENTIAL_TYPE_BASIC:  "basic",
	HTTP_CREDENTIAL_TYPE_BEARER: "bearer",
	HTTP_CREDENTIAL_TYPE_HEADER: "header",
	HTTP_CREDENTIAL_TYPE_QUERY:  "query",
}

var HTTP_CREDENTIAL_TYPE_CODES map[string]int

func init() {
	HTTP_CREDENTIAL_TYPE_CODES = make(map[string]int, len(HTTP_CREDENTIAL_TYPE_NAMES))
	for code, name := range HTTP_CREDENTIAL_TYPE_NAMES {
		HTTP_CREDENTIAL_TYPE_CODES[name] = code
	}
}

type HTTPCredential struct {
	Name  string
	Type  int
	Key   string
	Value string
}

var HTTP_CREDENTIALS Table[*HTTPCredential, HTTPCredential] = HTTPCredentialTable{
	TableHeader{
		Table:  "httpcredentials",
		IDs:    []string{"name"},
		Fields: []string{"type", "key", "value"},
	},
}

type HTTPCredentialTable struct {
	TableHeader
}

func (table HTTPCredentialTable) IDsValid(ids ...any) error {
	if len(ids) != 1 {
		return fmt.Errorf("invalid number of ids")
	}
	if name, ok := ids[0].(string); !ok || name == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (table HTTPCredentialTable) NewRecord() *HTTPCredential {
	return new(HTTPCredential)
}

func (record *HTTPCredential) IDPointers() []any {
	return []any{&record.Name}
}

func (record *HTTPCredential) FieldPointers() []any {
	return []any{&record.Type, &record.Key, &record.Value}
}

func (record *HTTPCredential) Dataset() HTTPCredential {
	return *record
}

func (record HTTPCredential) Valid() error {
	if record.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if _, typeOk := HTTP_CREDENTIAL_TYPE_NAMES[record.Type]; !typeOk {
		return fmt.Errorf("invalid type")
	}
	return nil
}

func (record HTTPCredential) IDs() []any {
	return record.IDPointers()
}

func (record HTTPCredential) Fields() []any {
	return record.FieldPointers()
}
