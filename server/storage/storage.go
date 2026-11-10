package storage

import (
	"github.com/2manyvcos/paranal/server/staticconfig"
	"github.com/2manyvcos/paranal/server/db"
	"github.com/2manyvcos/paranal/server/events"
)

type Storage struct {
	db             db.DB
	eventBroker    *events.EventBroker
	stateScheduler *stateScheduler
}

func (s *Storage) PutCollectorStates(collector string, states []State) error {
	return s.stateScheduler.PutCollectorStates(collector, states)
}

// CleanupCollectorStates removes stale states for collectors that do not exist anymore.
// Do not call this after starting any collector.
func (s *Storage) CleanupCollectorStates(collectors []string) error {
	_, err := s.db.CleanupCollectorStates(collectors)
	return err
}

func (s *Storage) Shutdown() {
	s.stateScheduler.Shutdown()
	s.db.Shutdown()
}

func New(cfg *staticconfig.Config, eventBroker *events.EventBroker) (*Storage, error) {
	var s Storage

	var err error
	s.db, err = db.NewSQLite(cfg.DB.File)
	if err != nil {
		return nil, err
	}

	s.stateScheduler = newStateScheduler(s.db, eventBroker)

	return &s, nil
}
