package application

import (
	"os"
	"strconv"
	"strings"

	"github.com/2manyvcos/paranal/utils"
)

const ENV_PREFIX = "PARANAL_"

func loadConfigValue(name string, defaultValue string) (result string) {
	result = os.Getenv(ENV_PREFIX + name)
	if result == "" {
		result = defaultValue
	}
	return
}

// func loadRequiredConfigValue(name string) (result string, err error) {
// 	result = loadConfigValue(name, "")
// 	if result == "" {
// 		err = fmt.Errorf("required environment variable %s is not set", ENV_PREFIX+name)
// 	}
// 	return
// }

func loadConfigBool(name string, defaultValue bool) (result bool) {
	result, _ = strconv.ParseBool(strings.ToLower(loadConfigValue(name, strconv.FormatBool(defaultValue))))
	return
}

func loadConfigList(name string, defaultValues []string) (result []string) {
	return utils.ParseList(loadConfigValue(name, ""), defaultValues)
}

// func loadConfigSecret(name string, defaultValue string) (result string, err error) {
// 	if value := loadConfigValue(name, ""); value != "" {
// 		return value, nil
// 	}

// 	if path := loadConfigValue(name+"_FILE", ""); path != "" {
// 		contents, err := os.ReadFile(path)
// 		return string(contents), err
// 	}

// 	return defaultValue, nil
// }
