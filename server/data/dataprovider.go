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
	UpdateUser(name string, record User) error
	DeleteUser(name string) error

	ListHTTPCredentials() ([]HTTPCredential, error)
	GetHTTPCredential(name string) (HTTPCredential, error)
	CreateHTTPCredential(record HTTPCredential, updateExisting bool) error
	UpdateHTTPCredential(name string, record HTTPCredential) error
	DeleteHTTPCredential(name string) error

	ListSSHCredentials() ([]SSHCredential, error)
	GetSSHCredential(name string) (SSHCredential, error)
	CreateSSHCredential(record SSHCredential, updateExisting bool) error
	UpdateSSHCredential(name string, record SSHCredential) error
	DeleteSSHCredential(name string) error

	ListUserCredentials() ([]UserCredential, error)
	GetUserCredential(name string) (UserCredential, error)
	CreateUserCredential(record UserCredential, updateExisting bool) error
	UpdateUserCredential(name string, record UserCredential) error
	DeleteUserCredential(name string) error

	GetService(id string) (Service, error)
	CreateService(record Service, updateExisting bool) error
	UpdateService(id string, record Service) error
	DeleteService(id string) error
	ListServicesWithFavorite(userName string) ([]ServiceWithFavorite, error)
	GetServiceWithFavorite(id string, userName string) (ServiceWithFavorite, error)

	CreateFavorite(record Favorite) error
	DeleteFavorite(userName string, serviceID string) error

	ListScripts() ([]Script, error)
	CreateScript(record Script, updateExisting bool) error
	ListScriptsByService(serviceID string) ([]Script, error)
	GetScriptByService(id int, serviceID string) (Script, error)
	UpdateScriptByService(id int, serviceID string, record Script) error
	DeleteScriptByService(id int, serviceID string) error
}
