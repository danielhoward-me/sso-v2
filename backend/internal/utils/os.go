package utils

import (
	"fmt"
	"os"
)

func IsDevelopment() bool {
	return GetEnv("ENV_NAME", "development") == "development"
}

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

func RequireEnv(key string) string {
	value, exists := os.LookupEnv(key)

	if exists {
		return value
	}

	panic(fmt.Errorf("%s is a required environment variable", key))
}
