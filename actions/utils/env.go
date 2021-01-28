package utils

import (
	"fmt"
	"os"
	"strconv"
)

func LoadEnvCritical(key string) string {
	if str, ok := os.LookupEnv(key); ok {
		return str
	} else {
		panic(fmt.Sprintf("Environment Variable \"%s\" is Required!", key))
	}
}

func LoadEnv(key string, defaultVal string) string {
	if str, ok := os.LookupEnv(key); ok {
		return str
	} else {
		return defaultVal
	}
}

func LoadEnvInt(key string, defaultVal int) int {
	if str, ok := os.LookupEnv(key); ok {
		if str, err := strconv.Atoi(str); err == nil {
			return str
		}
	}
	return defaultVal
}

func LoadEnvBool(key string, defaultVal bool) bool {
	if str, ok := os.LookupEnv(key); ok {
		if parseBool, err := strconv.ParseBool(str); err == nil {
			return parseBool
		}
	}
	return defaultVal
}
