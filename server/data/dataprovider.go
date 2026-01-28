package data

import "io"

type DataProvider interface {
	io.Closer

	InsertUser(user User) error
	UpsertUser(user User) error
	UpdateUser(user User) error
	GetUser(name string) (*User, error)
	ListUsers() ([]User, error)
}
