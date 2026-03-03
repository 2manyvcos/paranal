package schema

// index for services joined to service_configs

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
