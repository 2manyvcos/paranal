package data

import "fmt"

const USER_ROLE_COMMON = 1
const USER_ROLE_ADMIN = 2

var USER_ROLE_NAMES = map[int]string{
	USER_ROLE_COMMON: "common",
	USER_ROLE_ADMIN:  "admin",
}
var USER_ROLE_CODES map[string]int

func init() {
	USER_ROLE_CODES = make(map[string]int, len(USER_ROLE_NAMES))
	for code, name := range USER_ROLE_NAMES {
		USER_ROLE_CODES[name] = code
	}
}

type User struct {
	Name         string
	DisplayName  string
	Role         int
	PasswordHash string
}

var USERS Table[*User, User] = UserTable{
	TableHeader{
		Table:  "users",
		IDs:    []string{"name"},
		Fields: []string{"displayName", "role", "passwordHash"},
	},
}

type UserTable struct {
	TableHeader
}

func (table UserTable) IDsValid(ids ...any) error {
	if len(ids) != 1 {
		return fmt.Errorf("invalid number of ids")
	}
	if name, ok := ids[0].(string); !ok || name == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (table UserTable) NewRecord() *User {
	return new(User)
}

func (record *User) IDPointers() []any {
	return []any{&record.Name}
}

func (record *User) FieldPointers() []any {
	return []any{&record.DisplayName, &record.Role, &record.PasswordHash}
}

func (record *User) Dataset() User {
	return *record
}

func (record User) Valid() bool {
	_, roleOk := USER_ROLE_NAMES[record.Role]
	return record.Name != "" && roleOk
}

func (record User) IDs() []any {
	return record.IDPointers()
}

func (record User) Fields() []any {
	return record.FieldPointers()
}
