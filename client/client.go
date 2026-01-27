package client

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed all:dist
var rawClientFiles embed.FS

var ClientFiles fs.FS

func init() {
	var err error
	ClientFiles, err = fs.Sub(rawClientFiles, "dist")
	if err != nil {
		panic(fmt.Sprintf("forking client FS failed - %s\n", err))
	}
}
