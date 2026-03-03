package schema

// index for user_alert_channels joined to users and service_configs

type ServiceAlertChannelQuery struct {
	UserAlertChannelQuery
	UserQuery
	ServiceConfigQuery
	ErrorAlerts   *bool
	HealthAlerts  *bool
	VersionAlerts *bool
}

type ServiceAlertChannel struct {
	UserAlertChannel
}
