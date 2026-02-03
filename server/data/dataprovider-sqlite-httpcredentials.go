package data

import (
	"fmt"

	_ "modernc.org/sqlite"
)

func (p *sqliteImpl) setupHTTPCredentials() error {
	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS httpcredentials (id INTEGER PRIMARY KEY, name TEXT UNIQUE, type INTEGER, key TEXT, value TEXT)")
	if err != nil {
		return fmt.Errorf("creating table \"httpcredentials\" failed - %s", err)
	}
	return nil
}

func (p *sqliteImpl) ListHTTPCredentials() ([]HTTPCredential, error) {
	return sqliteListDatasets(p, HTTPCredentials)
}

func (p *sqliteImpl) GetHTTPCredential(name string) (result HTTPCredential, err error) {
	return sqliteGetDataset(p, HTTPCredentials, HTTPCredentialName(name))
}

func (p *sqliteImpl) CreateHTTPCredential(record HTTPCredential, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, HTTPCredentials, record, HTTPCredentialNameIDs)
	}
	return sqliteCreateDataset(p, HTTPCredentials, record)
}

func (p *sqliteImpl) UpdateHTTPCredential(record HTTPCredential) error {
	return sqliteUpdateDataset(p, HTTPCredentials, HTTPCredentialName(record.Name), record)
}

func (p *sqliteImpl) DeleteHTTPCredential(name string) error {
	return sqliteDeleteDataset(p, HTTPCredentials, HTTPCredentialName(name))
}
