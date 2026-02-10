package schema

// index for useralertchannels joined to users, services and serviceuserconfigs

type ServiceUserAlertChannelQuery struct {
	ErrorAlerts          *bool
	UptimeAlerts         *bool
	VersionAlerts        *bool
	UserRole             *int
	UserErrorAlerts      *bool
	UserUptimeAlerts     *bool
	UserVersionAlerts    *bool
	ServiceHidden        *bool
	ServiceUptimeAlerts  *bool
	ServiceVersionAlerts *bool
}

type ServiceUserAlertChannel struct {
	UserAlertChannel
}
