package db

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/2manyvcos/paranal/utils"
	_ "modernc.org/sqlite"
)

var _ DB = (*sqliteDB)(nil)

type sqliteDB struct{ *sql.DB }

func (db *sqliteDB) PutCollectorStates(collector string, states []State) (modifiedStates []StateID, err error) {
	if collector == "" {
		return nil, fmt.Errorf("missing collector")
	}

	return nil, fmt.Errorf("TODO:")

	// 	// TODO: build db.State bundle with revs

	// 	revs, err := u.db.GetStateRevs(&db.StateQuery{Collector: &collector})
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// TODO: remove matching revs from both the state bundle and revs

	// 	fmt.Println(collector, states) // TODO: remove

	// 	/*
	// 	   TODO:
	// 	   - single goroutine for updating states
	// 	   - validation happens outside
	// 	   - error / success is notified back and then handled outside
	// 	   - the following happens inside:
	// 	     - hashing (rev building)
	// 	     - updating existing entries and creating new entries (in one transaction)
	// 	     - deleting old entries
	// 	     - sending events (in batch)
	// 	   - first, a rev bundle for all IDs for the collector are fetched
	// 	   - unchanged revs are removed from both the update batch as well as the rev bundle
	// 	   - all remaining entries in the update batch are upserted into the database
	// 	   - all entries in the rev bundle that are missing in the update batch are removed from the database
	// 	   - change events for all entries in the update batch + rev bundle are send
	//     - revs may also be used outside for E-Tags if feasible (otherwise, the rev logic could be moved to db)
	// 	*/

	// 	return nil
	// }
}

func (db *sqliteDB) Shutdown() {
	db.Close()
}

func (db *sqliteDB) setup() error {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS states (
      registry TEXT NOT NULL,
      id TEXT NOT NULL,
      rev TEXT NOT NULL,
      collector TEXT NOT NULL,
      data TEXT NOT NULL,
      PRIMARY KEY (registry, id)
    )
      WITHOUT ROWID;
  `)
	return err
}

func NewSQLite(path string) (DB, error) {
	log.Printf("configuring SQLite database at %q\n", path)

	err := utils.MkdirConfig(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	var db sqliteDB

	db.DB, err = sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening SQLite database: %v", err)
	}

	err = db.setup()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("setting up SQLite database: %v", err)
	}

	return &db, nil
}
