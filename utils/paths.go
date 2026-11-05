package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func ConfigDir(appName string) (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return "", fmt.Errorf("locating config directory: %w", errors.Join(err, homeErr))
		}
		base = filepath.Join(home, ".config")
	}

	dir := filepath.Join(base, appName)

	if !filepath.IsAbs(dir) {
		return "", fmt.Errorf("config directory %q is not an absolute path", dir)
	}

	return dir, nil
}

func MkdirConfig(path string) error {
	if err := os.MkdirAll(path, 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	return nil
}
