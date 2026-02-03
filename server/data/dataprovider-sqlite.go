package data

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type sqliteImpl struct {
	*sql.DB
}

func SQLite(path string) (DataProvider, error) {
	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	provider := &sqliteImpl{db}
	return provider, provider.Setup()
}

func (p *sqliteImpl) Setup() error {
	log.Println("Setting up SQLite data backend")

	// Migration strategy as taken from a stackoverflow comment:

	// > For tracking the database version, I use the built in user-version variable that sqlite provides (sqlite does nothing with this variable, you are free to use it however you please). It starts at 0, and you can get/set this variable with the following sqlite statements:
	// >
	// > > PRAGMA user_version;
	// > > PRAGMA user_version = 1;
	// >
	// > When the app starts, I check the current user-version, apply any changes that are needed to bring the schema up to date, and then update the user-version. I wrap the updates in a transaction so that if anything goes wrong, the changes aren't committed.
	// >
	// > For making schema changes, sqlite supports "ALTER TABLE" syntax for certain operations (renaming the table or adding a column). This is an easy way to update existing tables in-place. See the documentation here: http://www.sqlite.org/lang_altertable.html. For deleting columns or other changes that aren't supported by the "ALTER TABLE" syntax, I create a new table, migrate data into it, drop the old table, and rename the new table to the original name.

	// var dbVersion int
	// err = p.DB.QueryRow("PRAGMA user_version").Scan(&dbVersion)

	if err := p.setupUsers(); err != nil {
		return err
	}

	if err := p.setupHTTPCredentials(); err != nil {
		return err
	}

	if err := p.setupSSHCredentials(); err != nil {
		return err
	}

	if err := p.setupUserCredentials(); err != nil {
		return err
	}

	if err := p.setupServices(); err != nil {
		return err
	}

	if err := p.setupScripts(); err != nil {
		return err
	}

	return nil
}

func requireFound(result sql.Result) error {
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}
