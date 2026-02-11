package schema

// index for useralertchannels joined to users and serviceuserconfigs

type ServiceUserAlertChannelQuery struct {
	UserAlertChannelQuery
	UserQuery
	ServiceUserConfigQuery
	ErrorAlerts   *bool
	UptimeAlerts  *bool
	VersionAlerts *bool
}

type ServiceUserAlertChannel struct {
	UserAlertChannel
}
