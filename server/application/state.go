package application

import "github.com/2manyvcos/paranal/server/data/schema"

type State interface {
	OnServiceDeleted(serviceID string)
	OnServiceScriptChanged(script schema.ServiceScript)
	OnServiceScriptDeleted(serviceId string, scriptID string)
}
