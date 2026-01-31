package data

import "fmt"

type UserCredential struct {
	Name        string
	Description string
	Value       string
}

var USER_CREDENTIALS Table[*UserCredential, UserCredential] = UserCredentialTable{
	TableHeader{
		Table:  "usercredentials",
		IDs:    []string{"name"},
		Fields: []string{"description", "value"},
	},
}

type UserCredentialTable struct {
	TableHeader
}

func (table UserCredentialTable) IDsValid(ids ...any) error {
	if len(ids) != 1 {
		return fmt.Errorf("invalid number of ids")
	}
	if name, ok := ids[0].(string); !ok || name == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (table UserCredentialTable) NewRecord() *UserCredential {
	return new(UserCredential)
}

func (record *UserCredential) IDPointers() []any {
	return []any{&record.Name}
}

func (record *UserCredential) FieldPointers() []any {
	return []any{&record.Description, &record.Value}
}

func (record *UserCredential) Dataset() UserCredential {
	return *record
}

func (record UserCredential) Valid() error {
	if record.Name == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (record UserCredential) IDs() []any {
	return record.IDPointers()
}

func (record UserCredential) Fields() []any {
	return record.FieldPointers()
}
