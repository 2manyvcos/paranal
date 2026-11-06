package db

type StateID struct {
	Registry string
	ID       string
}

type State struct {
	Registry  string
	ID        string
	Rev       string
	Collector string
	Data      []byte
}
