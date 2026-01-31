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
	return sqliteListDatasets(p, SSH_CREDENTIALS)
}

func (p *sqliteImpl) GetSSHCredential(name string) (result SSHCredential, err error) {
	return sqliteGetDataset(p, SSH_CREDENTIALS, name)
}

func (p *sqliteImpl) CreateSSHCredential(record SSHCredential, updateExisting bool) error {
	return sqliteCreateDataset(p, SSH_CREDENTIALS, record, updateExisting)
}

func (p *sqliteImpl) UpdateSSHCredential(record SSHCredential) error {
	return sqliteUpdateDataset(p, SSH_CREDENTIALS, record)
}

func (p *sqliteImpl) DeleteSSHCredential(name string) error {
	return sqliteDeleteDataset(p, SSH_CREDENTIALS, name)
}
