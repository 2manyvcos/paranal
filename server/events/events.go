package events

type Event interface {
	Event() string
	Data() string
}

var _ Event = (*ChangeEvent)(nil)

type ChangeEvent struct {
	path string
}

func (e ChangeEvent) Event() string {
	return "change"
}

func (e ChangeEvent) Data() string {
	return e.path
}

func NewChangeEvent(path string) Event {
	return ChangeEvent{
		path: path,
	}
}
