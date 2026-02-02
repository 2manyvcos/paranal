package utils

import "regexp"

var JsonRegex = regexp.MustCompile("^application/[^+]*[+]?(json);?.*$")
