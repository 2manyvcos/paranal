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

func (db *sqliteDB) Shutdown() {
	db.Close()
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
