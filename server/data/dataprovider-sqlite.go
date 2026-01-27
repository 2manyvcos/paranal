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

	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT UNIQUE, role INTEGER, pwHash TEXT)")
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
		"INSERT INTO users (name, role, pwHash) VALUES(?, ?, ?) ON CONFLICT(name) DO UPDATE SET role=excluded.role, pwHash=excluded.pwHash",
		user.Name, user.Role, user.PasswordHash,
	)
	return err
}

func (p *sqlite) ListUsers() ([]UserDataset, error) {
	rows, err := p.DB.Query("SELECT id, name, role, pwHash FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []UserDataset

	for rows.Next() {
		var user UserDataset
		if err := rows.Scan(&user.ID, &user.Name, &user.Role, &user.PasswordHash); err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
