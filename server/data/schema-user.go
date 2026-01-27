package data

const USER_ROLE_ADMIN = 0
const USER_ROLE_COMMON = 1

type User struct {
	Name         string
	Role         int
	PasswordHash string
}

func (u User) Valid() bool {
	return u.Name != "" && u.Role >= 0 && u.Role <= 1
}

type UserDataset struct {
	ID int
	User
}
