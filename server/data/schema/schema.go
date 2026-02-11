package schema

import (
	"errors"
	"io"
)

var ErrConflict = errors.New("conflicting record found")
var ErrNotFound = errors.New("no record found")

type DataProvider interface {
	io.Closer

	GetSetting(name string) (string, error)
	CreateOrUpdateSetting(name string, value string) error

	ListUsers(query *UserQuery) ([]User, error)
	GetUser(query UserQuery) (User, error)
	CreateUser(record User) error
	CreateOrUpdateUser(record User) error
	UpdateUsers(query UserQuery, record User) error
	DeleteUsers(query UserQuery) error

	ListUserAlertChannels(query *UserAlertChannelQuery) ([]UserAlertChannel, error)
	GetUserAlertChannel(query UserAlertChannelQuery) (UserAlertChannel, error)
	CreateUserAlertChannel(record UserAlertChannel) error
	UpdateUserAlertChannels(query UserAlertChannelQuery, record UserAlertChannel) error
	DeleteUserAlertChannels(query UserAlertChannelQuery) error

	ListHTTPCredentials(query *HTTPCredentialQuery) ([]HTTPCredential, error)
	GetHTTPCredential(query HTTPCredentialQuery) (HTTPCredential, error)
	CreateHTTPCredential(record HTTPCredential) error
	UpdateHTTPCredentials(query HTTPCredentialQuery, record HTTPCredential) error
	DeleteHTTPCredentials(query HTTPCredentialQuery) error

	ListSSHCredentials(query *SSHCredentialQuery) ([]SSHCredential, error)
	GetSSHCredential(query SSHCredentialQuery) (SSHCredential, error)
	CreateSSHCredential(record SSHCredential) error
	UpdateSSHCredentials(query SSHCredentialQuery, record SSHCredential) error
	DeleteSSHCredentials(query SSHCredentialQuery) error

	ListUserCredentials(query *UserCredentialQuery) ([]UserCredential, error)
	GetUserCredential(query UserCredentialQuery) (UserCredential, error)
	CreateUserCredential(record UserCredential) error
	UpdateUserCredentials(query UserCredentialQuery, record UserCredential) error
	DeleteUserCredentials(query UserCredentialQuery) error

	ListServices(query *ServiceQuery) ([]Service, error)
	GetService(query ServiceQuery) (Service, error)
	CreateService(record Service) error
	UpdateServices(query ServiceQuery, record Service) error
	DeleteServices(query ServiceQuery) error

	ListServiceConfigs(query *ServiceConfigQuery) ([]ServiceConfig, error)
	CreateOrUpdateServiceConfig(record ServiceConfig) error

	ListUserServices(userName string, query *UserServiceQuery) ([]UserService, error)
	GetUserService(userName string, query UserServiceQuery) (UserService, error)

	ListServiceAlertChannels(serviceID string, query *ServiceAlertChannelQuery) ([]ServiceAlertChannel, error)

	ListServiceScripts(query *ServiceScriptQuery) ([]ServiceScript, error)
	GetServiceScript(query ServiceScriptQuery) (ServiceScript, error)
	CreateServiceScript(record ServiceScript) error
	UpdateServiceScripts(query ServiceScriptQuery, record ServiceScript) error
	DeleteServiceScripts(query ServiceScriptQuery) error
}
