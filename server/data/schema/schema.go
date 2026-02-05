package schema

import (
	"errors"
	"io"
)

var ErrConflict = errors.New("conflicting record found")
var ErrNotFound = errors.New("no record found")

type DataProvider interface {
	io.Closer

	ListUsers(query *UserQuery) ([]User, error)
	GetUser(query UserQuery) (User, error)
	CreateUser(record User) error
	CreateOrUpdateUser(record User) error
	UpdateUser(query UserQuery, record User) error
	DeleteUser(query UserQuery) error

	ListHTTPCredentials(query *HTTPCredentialQuery) ([]HTTPCredential, error)
	GetHTTPCredential(query HTTPCredentialQuery) (HTTPCredential, error)
	CreateHTTPCredential(record HTTPCredential) error
	UpdateHTTPCredential(query HTTPCredentialQuery, record HTTPCredential) error
	DeleteHTTPCredential(query HTTPCredentialQuery) error

	ListSSHCredentials(query *SSHCredentialQuery) ([]SSHCredential, error)
	GetSSHCredential(query SSHCredentialQuery) (SSHCredential, error)
	CreateSSHCredential(record SSHCredential) error
	UpdateSSHCredential(query SSHCredentialQuery, record SSHCredential) error
	DeleteSSHCredential(query SSHCredentialQuery) error

	ListUserCredentials(query *UserCredentialQuery) ([]UserCredential, error)
	GetUserCredential(query UserCredentialQuery) (UserCredential, error)
	CreateUserCredential(record UserCredential) error
	UpdateUserCredential(query UserCredentialQuery, record UserCredential) error
	DeleteUserCredential(query UserCredentialQuery) error

	ListServicesWithFavorite(userName string, query *ServiceQuery) ([]ServiceWithFavorite, error)
	GetService(query ServiceQuery) (Service, error)
	GetServiceWithFavorite(userName string, query ServiceQuery) (ServiceWithFavorite, error)
	CreateService(record Service) error
	UpdateService(query ServiceQuery, record Service) error
	DeleteService(query ServiceQuery) error

	CreateFavorite(record Favorite) error
	DeleteFavorite(query FavoriteQuery) error

	ListScripts(query *ScriptQuery) ([]Script, error)
	GetScript(query ScriptQuery) (Script, error)
	CreateScript(record Script) error
	UpdateScript(query ScriptQuery, record Script) error
	DeleteScript(query ScriptQuery) error
}
