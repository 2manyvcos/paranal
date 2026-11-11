package storage

import (
	"sync"

	"github.com/2manyvcos/paranal/server/db"
	"github.com/2manyvcos/paranal/server/staticconfig"
)

type Storage struct {
	db        db.DB
	stateLock sync.Mutex
}

func (s *Storage) Shutdown() {
	s.db.Shutdown()
}

func New(cfg *staticconfig.Config) (*Storage, error) {
	var s Storage

	var err error
	s.db, err = db.NewSQLite(cfg.DB.File)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
