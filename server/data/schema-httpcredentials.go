package data

import "fmt"

type HTTPCredential struct {
	Name  string
	Key   string
	Value string
}

var HTTP_CREDENTIALS Table[*HTTPCredential, HTTPCredential] = HTTPCredentialTable{
	TableHeader{
		Table:  "httpcredentials",
		IDs:    []string{"name"},
		Fields: []string{"key", "value"},
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
	return []any{&record.Key, &record.Value}
}

func (record *HTTPCredential) Dataset() HTTPCredential {
	return *record
}

func (record HTTPCredential) Valid() bool {
	return record.Name != ""
}

func (record HTTPCredential) IDs() []any {
	return record.IDPointers()
}

func (record HTTPCredential) Fields() []any {
	return record.FieldPointers()
}
