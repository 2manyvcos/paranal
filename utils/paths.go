package utils

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
)

func ConfigDir(appName string) (string, error) {
	dir := filepath.Join(xdg.ConfigHome, appName)

	if !filepath.IsAbs(dir) {
		return "", fmt.Errorf("%q is not an absolute path", dir)
	}

	return dir, nil
}

func DataDir(appName string) (string, error) {
	dir := filepath.Join(xdg.DataHome, appName)

	if !filepath.IsAbs(dir) {
		return "", fmt.Errorf("%q is not an absolute path", dir)
	}

	return dir, nil
}
