package data

import (
	"fmt"

	"github.com/2manyvcos/paranal/server/data/sqlite"
	"github.com/2manyvcos/paranal/server/schema"
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
