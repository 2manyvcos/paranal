package data

import (
	"fmt"

	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/2manyvcos/paranal/server/data/sqlite"
)

type Config struct {
	Type string
	Path string
}

func Load(config Config) (schema.DataProvider, error) {
	switch config.Type {
	case "sqlite":
		return sqlite.NewProvider(config.Path)

	default:
		return nil, fmt.Errorf("invalid DB type \"%s\"", config.Type)
	}
}
