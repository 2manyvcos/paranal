package schema

// index for services joined to serviceuserconfigs

type UserServiceQuery struct {
	ID *string
}

type UserService struct {
	Service
	Favorite      bool
	Hidden        bool
	UptimeAlerts  bool
	VersionAlerts bool
}
