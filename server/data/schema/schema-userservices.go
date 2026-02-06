package schema

type UserServiceQuery struct {
	ID *string
}

type UserService struct {
	Service
	Favorite     bool
	Hidden       bool
	UptimeAlert  bool
	VersionAlert bool
}
