package data

import "fmt"

type SSHCredentialTable struct{ TableHeader }

var SSHCredentials Table[SSHCredential] = SSHCredentialTable{
	TableHeader{
		Name:   "sshcredentials",
		Fields: sshCredentialFields,
	},
}

func (t SSHCredentialTable) NewRecord() TableRecord[SSHCredential] {
	return new(SSHCredential)
}

type SSHCredential struct {
	Name       string
	User       string
	Password   string
	PrivateKey string
}

var sshCredentialFields = []string{"name", "user", "password", "privateKey"}

func (r *SSHCredential) RecordFields() []any {
	return []any{&r.Name, &r.User, &r.Password, &r.PrivateKey}
}

func (r *SSHCredential) Dataset() SSHCredential {
	return *r
}

var _ Upsertable[SSHCredential] = SSHCredential{}

func (r SSHCredential) Table() Table[SSHCredential] {
	return SSHCredentials
}

func (r SSHCredential) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if r.User == "" {
		return fmt.Errorf("invalid user")
	}
	if r.Password == "" && r.PrivateKey == "" {
		return fmt.Errorf("invalid credential")
	}
	return nil
}

func (r SSHCredential) InsertableNames() []string {
	return sshCredentialInsertables
}

var sshCredentialInsertables = []string{"name", "user", "password", "privateKey"}

func (r SSHCredential) Insertables() []any {
	return []any{r.Name, r.User, r.Password, r.PrivateKey}
}

func (r SSHCredential) UpdatableNames() []string {
	return sshCredentialUpdatables
}

var sshCredentialUpdatables = []string{"name", "user", "password", "privateKey"}

func (r SSHCredential) Updatables() []any {
	return []any{r.Name, r.User, r.Password, r.PrivateKey}
}

type SSHCredentialName string

var _ Identifier[SSHCredential] = SSHCredentialName("")

func (i SSHCredentialName) Table() Table[SSHCredential] {
	return SSHCredentials
}

func (i SSHCredentialName) Valid() error {
	if i == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (i SSHCredentialName) IDNames() []string {
	return SSHCredentialNameIDs
}

var SSHCredentialNameIDs = []string{"name"}

func (i SSHCredentialName) IDs() []any {
	return []any{string(i)}
}
