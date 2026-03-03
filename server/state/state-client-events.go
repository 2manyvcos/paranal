package state

import (
	"context"

	"github.com/2manyvcos/paranal/server/schema"
)

func (s *State) PublishClientEvent(event schema.ClientEvent) {
	s.clientEvents.Publish(event)
}

func (s *State) SubscribeToClientEvents(ctx context.Context) <-chan schema.ClientEvent {
	sub := s.clientEvents.Subscribe()
	if ctx != nil {
		go func() {
			<-ctx.Done()
			s.clientEvents.Unsubscribe(sub)
		}()
	}
	return sub
}
