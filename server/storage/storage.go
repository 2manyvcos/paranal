package storage

import (
	"github.com/2manyvcos/paranal/server/config"
	"github.com/2manyvcos/paranal/server/db"
)

type Storage struct {
	db db.DB
}

func (s *Storage) Shutdown() {
	s.db.Shutdown()
}

func New(cfg *config.Config) (*Storage, error) {
	var s Storage

	var err error
	s.db, err = db.NewSQLite(cfg.DB.File)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

/*
TODO:
- single goroutine for updating states
- validation happens outside
- error / success is notified back and then handled outside
- the following happens inside:
  - hashing (rev building)
  - updating existing entries and creating new entries
  - deleting old entries
  - sending events (in batch)
- first, a rev bundle for all IDs for the collector are fetched
- unchanged revs are removed from both the update batch as well as the rev bundle
- all remaining entries in the update batch are upserted into the database
- all entries in the rev bundle that are missing in the update batch are removed from the database
- change events for all entries in the update batch + rev bundle are send
*/
