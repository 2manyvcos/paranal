package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type sqlite struct {
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

	provider := &sqlite{db}
	return provider, provider.Setup()
}

func (p *sqlite) Setup() error {
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

	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT UNIQUE, displayName TEXT, role INTEGER, pwHash TEXT)")
	if err != nil {
		return fmt.Errorf(`creating table "users" failed - %s`, err)
	}

	return nil
}

func (p *sqlite) UpsertUser(user User) error {
	if !user.Valid() {
		return fmt.Errorf(`invalid user`)
	}

	_, err := p.DB.Exec(
		"INSERT INTO users (name, displayName, role, pwHash) VALUES(?, ?, ?, ?) ON CONFLICT(name) DO UPDATE SET displayName=excluded.displayName, role=excluded.role, pwHash=excluded.pwHash",
		user.Name, user.DisplayName, user.Role, user.PasswordHash,
	)
	return err
}

func (p *sqlite) ListUsers() ([]UserDataset, error) {
	rows, err := p.DB.Query("SELECT id, name, displayName, role, pwHash FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []UserDataset

	for rows.Next() {
		var user UserDataset
		if err := rows.Scan(&user.ID, &user.Name, &user.DisplayName, &user.Role, &user.PasswordHash); err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
