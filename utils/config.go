package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const EnvPrefix = "PARANAL_"

func LoadConfigValue(name string, defaultValue string) (result string) {
	result = os.Getenv(EnvPrefix + name)
	if result == "" {
		result = defaultValue
	}
	return
}

func LoadRequiredConfigValue(name string) (result string, err error) {
	result = LoadConfigValue(name, "")
	if result == "" {
		err = fmt.Errorf("required environment variable %s is not set", EnvPrefix+name)
	}
	return
}

func LoadConfigBool(name string, defaultValue bool) (result bool) {
	result, _ = strconv.ParseBool(strings.ToLower(LoadConfigValue(name, strconv.FormatBool(defaultValue))))
	return
}

func LoadConfigList(name string, defaultValues []string) (result []string) {
	return ParseList(LoadConfigValue(name, ""), defaultValues)
}

// func LoadConfigSecret(name string, defaultValue string) (result string, err error) {
// 	if value := LoadConfigValue(name, ""); value != "" {
// 		return value, nil
// 	}

// 	if path := LoadConfigValue(name+"_FILE", ""); path != "" {
// 		contents, err := os.ReadFile(path)
// 		return string(contents), err
// 	}

// 	return defaultValue, nil
// }
