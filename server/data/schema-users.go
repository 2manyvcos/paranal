package data

import "fmt"

const (
	USER_ROLE_COMMON = 1
	USER_ROLE_ADMIN  = 2
)

var (
	USER_ROLE_NAMES = map[int]string{
		USER_ROLE_COMMON: "common",
		USER_ROLE_ADMIN:  "admin",
	}
	USER_ROLE_CODES map[string]int
)

func init() {
	USER_ROLE_CODES = make(map[string]int, len(USER_ROLE_NAMES))
	for code, name := range USER_ROLE_NAMES {
		USER_ROLE_CODES[name] = code
	}
}

type UserTable struct{ TableHeader }

var Users Table[User] = UserTable{
	TableHeader{
		Name:   "users",
		Fields: userFields,
	},
}

func (t UserTable) NewRecord() TableRecord[User] {
	return new(User)
}

type User struct {
	Name         string
	DisplayName  string
	Role         int
	PasswordHash string
}

var userFields = []string{"name", "displayName", "role", "passwordHash"}

func (r *User) RecordFields() []any {
	return []any{&r.Name, &r.DisplayName, &r.Role, &r.PasswordHash}
}

func (r *User) Dataset() User {
	return *r
}

var _ Upsertable[User] = User{}

func (r User) Table() Table[User] {
	return Users
}

func (r User) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if _, roleOk := USER_ROLE_NAMES[r.Role]; !roleOk {
		return fmt.Errorf("invalid role")
	}
	return nil
}

func (r User) InsertableNames() []string {
	return userInsertables
}

var userInsertables = []string{"name", "displayName", "role", "passwordHash"}

func (r User) Insertables() []any {
	return []any{r.Name, r.DisplayName, r.Role, r.PasswordHash}
}

func (r User) UpdatableNames() []string {
	return userUpdatables
}

var userUpdatables = []string{"displayName", "role", "passwordHash"}

func (r User) Updatables() []any {
	return []any{r.DisplayName, r.Role, r.PasswordHash}
}

type UserName string

var _ Identifier[User] = UserName("")

func (i UserName) Table() Table[User] {
	return Users
}

func (i UserName) Valid() error {
	if i == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (i UserName) IDNames() []string {
	return UserNameIDs
}

var UserNameIDs = []string{"name"}

func (i UserName) IDs() []any {
	return []any{string(i)}
}
