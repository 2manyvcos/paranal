package data

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

func (u User) Valid() bool {
	_, roleOk := USER_ROLE_NAMES[u.Role]
	return u.Name != "" && roleOk
}
