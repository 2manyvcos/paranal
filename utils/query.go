package utils

import (
	"net/url"
	"strconv"
)

func LoadQueryValue(query url.Values, name string) (string, bool) {
	if !query.Has(name) {
		return "", false
	}
	return query.Get(name), true
}

func LoadQueryBool(query url.Values, name string) (bool, bool) {
	v, ok := LoadQueryValue(query, name)
	if !ok {
		return false, false
	}
	if v == "" {
		return true, true
	}
	result, _ := strconv.ParseBool(v)
	return result, true
}
