package data

import (
	"fmt"

	_ "modernc.org/sqlite"
)

func (p *sqliteImpl) setupUsers() error {
	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT UNIQUE, displayName TEXT, role INTEGER, passwordHash TEXT)")
	if err != nil {
		return fmt.Errorf("creating table \"users\" failed - %s", err)
	}
	return nil
}

func (p *sqliteImpl) ListUsers() ([]User, error) {
	return sqliteListDatasets(p, USERS)
}

func (p *sqliteImpl) GetUser(name string) (result User, err error) {
	return sqliteGetDataset(p, USERS, name)
}

func (p *sqliteImpl) CreateUser(record User, updateExisting bool) error {
	return sqliteCreateDataset(p, USERS, record, updateExisting)
}

func (p *sqliteImpl) UpdateUser(record User) error {
	return sqliteUpdateDataset(p, USERS, record)
}

func (p *sqliteImpl) DeleteUser(name string) error {
	return sqliteDeleteDataset(p, USERS, name)
}
