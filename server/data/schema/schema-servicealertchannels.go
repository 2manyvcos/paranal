package schema

// index for useralertchannels joined to users and serviceconfigs

type ServiceAlertChannelQuery struct {
	UserAlertChannelQuery
	UserQuery
	ServiceConfigQuery
	ErrorAlerts   *bool
	UptimeAlerts  *bool
	VersionAlerts *bool
}

type ServiceAlertChannel struct {
	UserAlertChannel
}
