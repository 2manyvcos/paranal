package schema

func NewUpdateEvent(path string) ClientEvent {
	return UpdateEvent{data: path}
}

func NewUserUpdateEvent(user string, path string) ClientEvent {
	return UpdateEvent{user: user, data: path}
}

type UpdateEvent struct {
	user string
	data string
}

func (e UpdateEvent) Event() string {
	return "update"
}

func (e UpdateEvent) User() string {
	return e.user
}

func (e UpdateEvent) Data() string {
	return e.data
}
