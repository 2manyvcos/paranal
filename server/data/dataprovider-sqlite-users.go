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
	return sqliteListDatasets(p, Users)
}

func (p *sqliteImpl) GetUser(name string) (result User, err error) {
	return sqliteGetDataset(p, Users, UserName(name))
}

func (p *sqliteImpl) CreateUser(record User, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, Users, record, UserNameIDs)
	}
	return sqliteCreateDataset(p, Users, record)
}

func (p *sqliteImpl) UpdateUser(record User) error {
	return sqliteUpdateDataset(p, Users, UserName(record.Name), record)
}

func (p *sqliteImpl) DeleteUser(name string) error {
	return sqliteDeleteDataset(p, Users, UserName(name))
}
