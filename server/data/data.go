package data

import (
	"fmt"
)

type Config struct {
	Type string
	Path string
}

func Load(config Config) (DataProvider, error) {
	switch config.Type {
	case "sqlite":
		return SQLite(config.Path)

	default:
		return nil, fmt.Errorf(`invalid DB type "%s"`, config.Type)
	}
}
