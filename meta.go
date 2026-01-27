package meta

import (
	_ "embed"
	"encoding/json"
	"log"
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
		log.Fatalf("Error reading package metadata - %s\n", err)
	}
}
