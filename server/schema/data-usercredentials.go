package schema

import (
	"fmt"
)

type UserCredentialQuery struct {
	Name        *string
	Description *string
	HasValue    *bool
}

type UserCredential struct {
	Name        string
	Description string
	Value       string
}

func (r UserCredential) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if r.Value == "" {
		return fmt.Errorf("invalid value")
	}
	return nil
}
