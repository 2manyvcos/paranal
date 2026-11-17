package dynamicconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/2manyvcos/paranal/spec"
	"go.yaml.in/yaml/v4"
	"golang.org/x/mod/semver"
)

type configEntry struct {
	Type    string         `yaml:"type"`
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
		config, err := parseConfig("file", entry)
		if err != nil {
			return nil, err
		}

		configs[i] = config
	}

	return configs, nil
}

func parseConfig(provider string, config configEntry) (Config, error) {
	if config.Type == "" {
		return nil, fmt.Errorf("missing \"type\"")
	}

	switch config.Type {
	case "collector":
		result, err := parseCollector(provider, config)
		if err != nil {
			return nil, fmt.Errorf("parsing collector: %w", err)
		}
		return result, nil

	case "dashboard":
		result, err := parseDashboard(provider, config)
		if err != nil {
			return nil, fmt.Errorf("parsing dashboard: %w", err)
		}
		return result, nil

	default:
		return nil, fmt.Errorf("unsupported type %q", config.Type)
	}
}

func parseVersion(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("missing \"version\"")
	}
	version := "v" + input
	if !semver.IsValid(version) || semver.MajorMinor(version) != version {
		return "", fmt.Errorf("invalid version %q", version)
	}
	return version, nil
}

func parseCollector(provider string, config configEntry) (Config, error) {
	version, err := parseVersion(config.Version)
	if err != nil {
		return nil, err
	}

	parser, err := spec.CollectionParser(version)
	if err != nil {
		return nil, err
	}

	result, err := parser.ParseCollector(provider, config.Config)
	if err != nil {
		return nil, err
	}

	return CollectorConfig{result}, nil
}

func parseDashboard(provider string, config configEntry) (Config, error) {
	version, err := parseVersion(config.Version)
	if err != nil {
		return nil, err
	}

	parser, err := spec.DashboardParser(version)
	if err != nil {
		return nil, err
	}

	result, err := parser.ParseDashboard(provider, config.Config)
	if err != nil {
		return nil, err
	}

	return DashboardConfig{result}, nil
}
