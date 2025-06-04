package pkg

import (
	"os"
	"strconv"
)

// GetEnv returns the value of the environment variable named by the key or return defaultValue
func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value

}

// GetEnvInt returns the value of the environment variable named by the key or return defaultValue
func GetEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue

}
