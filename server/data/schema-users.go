package data

import "fmt"

const (
	UserRoleCommon = iota + 1
	UserRoleAdmin
)

var (
	UserRoleNames = map[int]string{
		UserRoleCommon: "common",
		UserRoleAdmin:  "admin",
	}
	UserRoleCodes map[string]int
)

func init() {
	UserRoleCodes = make(map[string]int, len(UserRoleNames))
	for code, name := range UserRoleNames {
		UserRoleCodes[name] = code
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

func (r User) TableType() Table[User] { return nil }

func (r User) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if _, roleOk := UserRoleNames[r.Role]; !roleOk {
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
