package schema

import "fmt"

type HTTPCredentialQuery struct {
	Name *string
}

type HTTPCredential struct {
	Name  string
	Type  int
	Key   string
	Value string
}

const (
	HTTPCredentialTypeBasic = iota + 1
	HTTPCredentialTypeBearer
	HTTPCredentialTypeHeader
	HTTPCredentialTypeQuery
)

var (
	HTTPCredentialTypeNames = map[int]string{
		HTTPCredentialTypeBasic:  "basic",
		HTTPCredentialTypeBearer: "bearer",
		HTTPCredentialTypeHeader: "header",
		HTTPCredentialTypeQuery:  "query",
	}
	HTTPCredentialTypeCodes map[string]int
)

func init() {
	HTTPCredentialTypeCodes = make(map[string]int, len(HTTPCredentialTypeNames))
	for code, name := range HTTPCredentialTypeNames {
		HTTPCredentialTypeCodes[name] = code
	}
}

func (r HTTPCredential) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	switch r.Type {
	case HTTPCredentialTypeBasic:
		if r.Key == "" {
			return fmt.Errorf("invalid user")
		}
		if r.Value == "" {
			return fmt.Errorf("invalid password")
		}
	case HTTPCredentialTypeBearer:
		if r.Value == "" {
			return fmt.Errorf("invalid token")
		}
	case HTTPCredentialTypeHeader:
		if r.Key == "" {
			return fmt.Errorf("invalid header name")
		}
		if r.Value == "" {
			return fmt.Errorf("invalid header value")
		}
	case HTTPCredentialTypeQuery:
		if r.Key == "" {
			return fmt.Errorf("invalid query key")
		}
		if r.Value == "" {
			return fmt.Errorf("invalid query value")
		}
	default:
		return fmt.Errorf("invalid type")
	}
	return nil
}
