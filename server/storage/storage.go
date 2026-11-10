package storage

import (
	"sync"

	"github.com/2manyvcos/paranal/server/db"
	"github.com/2manyvcos/paranal/server/events"
	"github.com/2manyvcos/paranal/server/staticconfig"
)

type Storage struct {
	eventBroker *events.EventBroker
	db          db.DB
	stateLock   sync.Mutex
}

func (s *Storage) Shutdown() {
	s.db.Shutdown()
}

func New(cfg *staticconfig.Config, eventBroker *events.EventBroker) (*Storage, error) {
	var s Storage
	s.eventBroker = eventBroker

	var err error
	s.db, err = db.NewSQLite(cfg.DB.File)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
