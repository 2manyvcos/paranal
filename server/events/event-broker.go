package events

import (
	"github.com/2manyvcos/paranal/utils"
)

type EventBroker struct{ *utils.Broker[[]Event] }

func (b *EventBroker) Publish(e ...Event) {
	if len(e) > 0 {
		b.Broker.Publish(e)
	}
}

func NewEventBroker() *EventBroker {
	var b EventBroker
	b.Broker = utils.NewBroker[[]Event]()

	return &b
}
