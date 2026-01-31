package data

import (
	"fmt"

	_ "modernc.org/sqlite"
)

func (p *sqliteImpl) setupUserCredentials() error {
	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS usercredentials (id INTEGER PRIMARY KEY, name TEXT UNIQUE, description TEXT, value TEXT)")
	if err != nil {
		return fmt.Errorf("creating table \"usercredentials\" failed - %s", err)
	}
	return nil
}

func (p *sqliteImpl) ListUserCredentials() ([]UserCredential, error) {
	return sqliteListDatasets(p, USER_CREDENTIALS)
}

func (p *sqliteImpl) GetUserCredential(name string) (result UserCredential, err error) {
	return sqliteGetDataset(p, USER_CREDENTIALS, name)
}

func (p *sqliteImpl) CreateUserCredential(record UserCredential, updateExisting bool) error {
	return sqliteCreateDataset(p, USER_CREDENTIALS, record, updateExisting)
}

func (p *sqliteImpl) UpdateUserCredential(record UserCredential) error {
	return sqliteUpdateDataset(p, USER_CREDENTIALS, record)
}

func (p *sqliteImpl) DeleteUserCredential(name string) error {
	return sqliteDeleteDataset(p, USER_CREDENTIALS, name)
}
