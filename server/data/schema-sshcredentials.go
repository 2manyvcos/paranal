package data

import "fmt"

type SSHCredential struct {
	Name       string
	User       string
	Password   string
	PrivateKey string
}

var SSH_CREDENTIALS Table[*SSHCredential, SSHCredential] = SSHCredentialTable{
	TableHeader{
		Table:  "sshcredentials",
		IDs:    []string{"name"},
		Fields: []string{"user", "password", "privateKey"},
	},
}

type SSHCredentialTable struct {
	TableHeader
}

func (table SSHCredentialTable) IDsValid(ids ...any) error {
	if len(ids) != 1 {
		return fmt.Errorf("invalid number of ids")
	}
	if name, ok := ids[0].(string); !ok || name == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (table SSHCredentialTable) NewRecord() *SSHCredential {
	return new(SSHCredential)
}

func (record *SSHCredential) IDPointers() []any {
	return []any{&record.Name}
}

func (record *SSHCredential) FieldPointers() []any {
	return []any{&record.User, &record.Password, &record.PrivateKey}
}

func (record *SSHCredential) Dataset() SSHCredential {
	return *record
}

func (record SSHCredential) Valid() error {
  if record.Name == "" {
    return fmt.Errorf("invalid name")
  }
  if record.User == "" {
    return fmt.Errorf("invalid user")
  }
  return nil
}

func (record SSHCredential) IDs() []any {
	return record.IDPointers()
}

func (record SSHCredential) Fields() []any {
	return record.FieldPointers()
}
