package sqlite

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/2manyvcos/paranal/server/schema"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var setups []func(i *impl) error
var cleanups []func(i *impl) error

type impl struct{ *sql.DB }

func NewProvider(path string) (schema.DataProvider, error) {
	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	i := &impl{db}
	for _, setup := range setups {
		if err := setup(i); err != nil {
			return nil, err
		}
	}
	return i, nil
}

func requireNoConflict(err error) error {
	if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return schema.ErrConflict
	}
	return err
}

func requireFound(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return schema.ErrNotFound
	}
	return nil
}

func (i *impl) CleanupDatabase() {
	for _, cleanup := range cleanups {
		if err := cleanup(i); err != nil {
			log.Printf("Error cleaning up database - %s\n", err)
		}
	}
}
