package db

type err string

func (e err) Error() string {
	return string(e)
}

const (
	ErrNotFound = err("not found")
)
