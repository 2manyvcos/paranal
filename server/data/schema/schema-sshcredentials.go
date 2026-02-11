package schema

import "fmt"

type SSHCredentialQuery struct {
	Name *string
	User *string
}

type SSHCredential struct {
	Name       string
	User       string
	Password   string
	PrivateKey string
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
