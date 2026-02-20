package state

import (
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
)

func (s *State) setupClientEvents() error {
	s.clientEvents = utils.NewBroker[schema.ClientEvent]()
	go s.clientEvents.Start()

	return nil
}

func (s *State) closeClientEvents() {
	s.clientEvents.Close()
}
