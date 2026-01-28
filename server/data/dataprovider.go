package data

import (
	"errors"
	"io"
)

var ErrNotFound = errors.New("no record found")

type DataProvider interface {
	io.Closer

	ListUsers() ([]User, error)
	GetUser(name string) (User, error)
	CreateUser(user User, updateExisting bool) error
	UpdateUser(user User) error
	DeleteUser(name string) error
}
