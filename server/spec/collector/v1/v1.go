package collector_spec_v1

import collector_spec "github.com/2manyvcos/paranal/server/spec/collector"

const Version = "v1.0"

var Parser collector_spec.CollectorSpecParser = CollectorParserV1{}

type CollectorParserV1 struct{}
