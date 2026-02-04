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
	return sqliteSelectDatasets(p, Users, nil)
}

func (p *sqliteImpl) GetUser(name string) (result User, err error) {
	return sqliteSelectDataset(p, Users, &Conditions{Condition: &Condition{Field: "name", Value: name}})
}

func (p *sqliteImpl) CreateUser(record User, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, Users, record, []string{"name"})
	}
	return sqliteCreateDataset(p, Users, record)
}

func (p *sqliteImpl) UpdateUser(name string, record User) error {
	return sqliteUpdateDataset(p, Users, &Conditions{Condition: &Condition{Field: "name", Value: name}}, record)
}

func (p *sqliteImpl) DeleteUser(name string) error {
	return sqliteDeleteDataset(p, Users, &Conditions{Condition: &Condition{Field: "name", Value: name}})
}
