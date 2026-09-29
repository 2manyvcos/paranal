package utils

import "strings"

func ParseList(data string, defaultValues []string) (result []string) {
	if data == "" {
		return defaultValues
	}

	for _, part := range strings.Split(data, ",") {
		if value := strings.TrimSpace(part); value != "" {
			result = append(result, value)
		}
	}
	return
}
