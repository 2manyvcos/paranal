package data

import (
	"fmt"

	_ "modernc.org/sqlite"
)

func (p *sqliteImpl) setupHTTPCredentials() error {
	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS httpcredentials (id INTEGER PRIMARY KEY, name TEXT UNIQUE, key TEXT, value TEXT)")
	if err != nil {
		return fmt.Errorf("creating table \"httpcredentials\" failed - %s", err)
	}
	return nil
}

func (p *sqliteImpl) ListHTTPCredentials() ([]HTTPCredential, error) {
	return sqliteListDatasets(p, HTTP_CREDENTIALS)
}

func (p *sqliteImpl) GetHTTPCredential(name string) (result HTTPCredential, err error) {
	return sqliteGetDataset(p, HTTP_CREDENTIALS, name)
}

func (p *sqliteImpl) CreateHTTPCredential(record HTTPCredential, updateExisting bool) error {
	return sqliteCreateDataset(p, HTTP_CREDENTIALS, record, updateExisting)
}

func (p *sqliteImpl) UpdateHTTPCredential(record HTTPCredential) error {
	return sqliteUpdateDataset(p, HTTP_CREDENTIALS, record)
}

func (p *sqliteImpl) DeleteHTTPCredential(name string) error {
	return sqliteDeleteDataset(p, HTTP_CREDENTIALS, name)
}
