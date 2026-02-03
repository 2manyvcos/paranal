package data

import (
	"errors"
	"io"
)

var ErrConflict = errors.New("conflicting record found")
var ErrNotFound = errors.New("no record found")

type DataProvider interface {
	io.Closer

	ListUsers() ([]User, error)
	GetUser(name string) (User, error)
	CreateUser(record User, updateExisting bool) error
	UpdateUser(record User) error
	DeleteUser(name string) error

	ListHTTPCredentials() ([]HTTPCredential, error)
	GetHTTPCredential(name string) (HTTPCredential, error)
	CreateHTTPCredential(record HTTPCredential, updateExisting bool) error
	UpdateHTTPCredential(record HTTPCredential) error
	DeleteHTTPCredential(name string) error

	ListSSHCredentials() ([]SSHCredential, error)
	GetSSHCredential(name string) (SSHCredential, error)
	CreateSSHCredential(record SSHCredential, updateExisting bool) error
	UpdateSSHCredential(record SSHCredential) error
	DeleteSSHCredential(name string) error

	ListUserCredentials() ([]UserCredential, error)
	GetUserCredential(name string) (UserCredential, error)
	CreateUserCredential(record UserCredential, updateExisting bool) error
	UpdateUserCredential(record UserCredential) error
	DeleteUserCredential(name string) error

	ListServices() ([]Service, error)
	GetService(id string) (Service, error)
	CreateService(record Service, updateExisting bool) error
	UpdateService(record Service) error
	DeleteService(id string) error

	ListScripts() ([]Script, error)
	GetScript(name string) (Script, error)
	CreateScript(record Script, updateExisting bool) error
	UpdateScript(record Script) error
	DeleteScript(name string) error
}
