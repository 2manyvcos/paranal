package dynamicconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v4"
	"golang.org/x/mod/semver"
)

func loadConfigFile(path string) ([]Config, error) {
	var entries []configEntry
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = yaml.Load(
			contents,
			&entries,
			yaml.WithKnownFields(),
			yaml.WithUniqueKeys(),
			yaml.WithAllDocuments(),
		)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported format %q", ext)
	}

	configs := make([]Config, len(entries))
	for i, entry := range entries {
		config, err := parseConfig(entry)
		if err != nil {
			return nil, err
		}

		configs[i] = config
	}

	return configs, nil
}

type configEntry struct {
	Kind    string         `yaml:"kind"`
	Version string         `yaml:"version"`
	Config  map[string]any `yaml:",inline"`
}

func parseConfig(config configEntry) (Config, error) {
	switch config.Kind {
	case "collector":
		version := "v" + config.Version
		if !semver.IsValid(version) || semver.MajorMinor(version) != version {
			return nil, fmt.Errorf("invalid version %q", config.Version)
		}
		config, err := parseCollectorConfig(version, config.Config)
		if err != nil {
			return nil, fmt.Errorf("parsing collector: %w", err)
		}
		return config, nil

	default:
		return nil, fmt.Errorf("unsupported kind %q", config.Kind)
	}
}
