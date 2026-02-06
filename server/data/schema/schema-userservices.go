package schema

type UserServiceQuery struct {
	ID           *string
	Favorite     *bool
	UptimeAlert  *bool
	VersionAlert *bool
}

type UserService struct {
	Service
	Favorite     bool
	UptimeAlert  bool
	VersionAlert bool
}
