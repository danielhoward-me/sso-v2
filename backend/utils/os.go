package utils

import "os"

func GetEnv(key string, fallback ...string) string {
	value, exists := os.LookupEnv(key)

	if exists {
		return value
	} else if len(fallback) != 0 {
		return fallback[0]
	} else {
		return ""
	}
}
