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
	return sqliteListDatasets(p, SSHCredentials)
}

func (p *sqliteImpl) GetSSHCredential(name string) (result SSHCredential, err error) {
	return sqliteGetDataset(p, SSHCredentials, SSHCredentialName(name))
}

func (p *sqliteImpl) CreateSSHCredential(record SSHCredential, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, SSHCredentials, record, SSHCredentialNameIDs)
	}
	return sqliteCreateDataset(p, SSHCredentials, record)
}

func (p *sqliteImpl) UpdateSSHCredential(record SSHCredential) error {
	return sqliteUpdateDataset(p, SSHCredentials, SSHCredentialName(record.Name), record)
}

func (p *sqliteImpl) DeleteSSHCredential(name string) error {
	return sqliteDeleteDataset(p, SSHCredentials, SSHCredentialName(name))
}
