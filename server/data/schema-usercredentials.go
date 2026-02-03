package data

import "fmt"

type UserCredentialTable struct{ TableHeader }

var UserCredentials Table[UserCredential] = UserCredentialTable{
	TableHeader{
		Name:   "usercredentials",
		Fields: userCredentialFields,
	},
}

func (t UserCredentialTable) NewRecord() TableRecord[UserCredential] {
	return new(UserCredential)
}

type UserCredential struct {
	Name        string
	Description string
	Value       string
}

var userCredentialFields = []string{"name", "description", "value"}

func (r *UserCredential) RecordFields() []any {
	return []any{&r.Name, &r.Description, &r.Value}
}

func (r *UserCredential) Dataset() UserCredential {
	return *r
}

var _ Upsertable[UserCredential] = UserCredential{}

func (r UserCredential) Table() Table[UserCredential] {
	return UserCredentials
}

func (r UserCredential) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if r.Value == "" {
		return fmt.Errorf("invalid value")
	}
	return nil
}

func (r UserCredential) InsertableNames() []string {
	return userCredentialInsertables
}

var userCredentialInsertables = []string{"name", "description", "value"}

func (r UserCredential) Insertables() []any {
	return []any{r.Name, r.Description, r.Value}
}

func (r UserCredential) UpdatableNames() []string {
	return userCredentialUpdatables
}

var userCredentialUpdatables = []string{"name", "description", "value"}

func (r UserCredential) Updatables() []any {
	return []any{r.Name, r.Description, r.Value}
}

type UserCredentialName string

var _ Identifier[UserCredential] = UserCredentialName("")

func (i UserCredentialName) Table() Table[UserCredential] {
	return UserCredentials
}

func (i UserCredentialName) Valid() error {
	if i == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (i UserCredentialName) IDNames() []string {
	return UserCredentialNameIDs
}

var UserCredentialNameIDs = []string{"name"}

func (i UserCredentialName) IDs() []any {
	return []any{string(i)}
}
