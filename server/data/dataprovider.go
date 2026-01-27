package data

import "io"

type DataProvider interface {
	io.Closer

	UpsertUser(user User) error
	ListUsers() ([]UserDataset, error)
}
