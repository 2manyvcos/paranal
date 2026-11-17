package events

import "strings"

func MatchEvent(pattern, path string) bool {
	p := strings.Split(pattern, "/")
	e := strings.Split(path, "/")

	if len(p) != len(e) {
		return false
	}

	for i, ps := range p {
		es := e[i]
		if ps != "*" && ps != es {
			return false
		}
	}

	return true
}
