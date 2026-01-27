package meta

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed package.json
var rawPackage string

var Meta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func init() {
	err := json.Unmarshal([]byte(rawPackage), &Meta)
	if err != nil {
		panic(fmt.Sprintf("reading package metadata failed - %s\n", err))
	}
}
