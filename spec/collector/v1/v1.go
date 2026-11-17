package collector_spec_v1

import collector_spec "github.com/2manyvcos/paranal/spec/collector"

const Version = "v1.0"

var Parser collector_spec.CollectorSpecParser = collectorParser{}

type collectorParser struct{}
