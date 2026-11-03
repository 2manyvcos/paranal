package event

type Event interface {
	Event() string
	Data() string
}

var _ Event = (*ChangeEvent)(nil)

type ChangeEvent struct {
	path string
}

func (u ChangeEvent) Event() string {
	return "change"
}

func (u ChangeEvent) Data() string {
	return u.path
}

func NewChangeEvent(path string) Event {
	return ChangeEvent{
		path: path,
	}
}
