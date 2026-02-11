package schema

// index for useralertchannels joined to users and serviceconfigs

type ServiceUserAlertChannelQuery struct {
	UserAlertChannelQuery
	UserQuery
	ServiceConfigQuery
	ErrorAlerts   *bool
	UptimeAlerts  *bool
	VersionAlerts *bool
}

type ServiceUserAlertChannel struct {
	UserAlertChannel
}
