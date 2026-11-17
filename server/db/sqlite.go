package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var _ DB = (*sqliteDB)(nil)

type sqliteDB struct {
	*sql.DB
	sync.Mutex
}

func (db *sqliteDB) Shutdown() {
	db.Close()
}

func NewSQLite(path string) (DB, error) {
	log.Printf("configuring SQLite database at %q\n", path)

	err := os.MkdirAll(filepath.Dir(path), 0o700)
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
