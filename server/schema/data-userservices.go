package schema

// index for services joined to serviceconfigs

type UserServiceQuery struct {
	ServiceQuery
	ServiceConfigQuery
}

type UserService struct {
	Service
	Favorite      bool
	Hidden        bool
	HealthAlerts  bool
	VersionAlerts bool
}
