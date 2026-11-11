package spec

import (
	"fmt"

	collector_spec "github.com/2manyvcos/paranal/spec/collector"
	collector_spec_v1 "github.com/2manyvcos/paranal/spec/collector/v1"
	"golang.org/x/mod/semver"
)

func CollectionParser(version string) (collector_spec.CollectorSpecParser, error) {
	switch semver.Major(version) {
	case "v1":
		if semver.Compare(version, collector_spec_v1.Version) > 0 {
			return nil, fmt.Errorf("unsupported version %q", version)
		}
		return collector_spec_v1.Parser, nil

	default:
		return nil, fmt.Errorf("unsupported version %q", version)
	}
}
