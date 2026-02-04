package data

import (
	"fmt"

	_ "modernc.org/sqlite"
)

func (p *sqliteImpl) setupSSHCredentials() error {
	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS sshcredentials (id INTEGER PRIMARY KEY, name TEXT UNIQUE, user TEXT, password TEXT, privateKey TEXT)")
	if err != nil {
		return fmt.Errorf("creating table \"sshcredentials\" failed - %s", err)
	}
	return nil
}

func (p *sqliteImpl) ListSSHCredentials() ([]SSHCredential, error) {
	return sqliteSelectDatasets(p, SSHCredentials, nil)
}

func (p *sqliteImpl) GetSSHCredential(name string) (result SSHCredential, err error) {
	return sqliteSelectDataset(p, SSHCredentials, &Conditions{Condition: &Condition{Field: "name", Value: name}})
}

func (p *sqliteImpl) CreateSSHCredential(record SSHCredential, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, SSHCredentials, record, []string{"name"})
	}
	return sqliteCreateDataset(p, SSHCredentials, record)
}

func (p *sqliteImpl) UpdateSSHCredential(name string, record SSHCredential) error {
	return sqliteUpdateDataset(p, SSHCredentials, &Conditions{Condition: &Condition{Field: "name", Value: name}}, record)
}

func (p *sqliteImpl) DeleteSSHCredential(name string) error {
	return sqliteDeleteDataset(p, SSHCredentials, &Conditions{Condition: &Condition{Field: "name", Value: name}})
}
