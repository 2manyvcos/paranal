package dynamicconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/2manyvcos/paranal/server/spec"
	"go.yaml.in/yaml/v4"
	"golang.org/x/mod/semver"
)

type configEntry struct {
	Kind    string         `yaml:"kind"`
	Version string         `yaml:"version"`
	Config  map[string]any `yaml:",inline"`
}

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

func parseConfig(config configEntry) (Config, error) {
	switch config.Kind {
	case "collector":
		result, err := parseCollector(config)
		if err != nil {
			return nil, fmt.Errorf("parsing collector: %w", err)
		}
		return result, nil

	default:
		return nil, fmt.Errorf("unsupported kind %q", config.Kind)
	}
}

func parseCollector(config configEntry) (Config, error) {
	version := "v" + config.Version
	if !semver.IsValid(version) || semver.MajorMinor(version) != version {
		return nil, fmt.Errorf("invalid version %q", version)
	}

	parser, err := spec.CollectionParser(version)
	if err != nil {
		return nil, err
	}

	c, err := parser.ParseCollector(config.Config)
	if err != nil {
		return nil, err
	}

	return CollectorConfig(c), nil
}
