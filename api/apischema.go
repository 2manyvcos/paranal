package apischema

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed all:schema
var rawSchemaFiles embed.FS

var SchemaFiles fs.FS

func init() {
	var err error
	SchemaFiles, err = fs.Sub(rawSchemaFiles, "schema")
	if err != nil {
		panic(fmt.Sprintf("forking schema FS failed - %s\n", err))
	}
}
