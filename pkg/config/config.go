package config

import (
	"os"
	"strconv"
)

type conf struct {
	PORT             uint
	JWT_TOKEN        string
	DB_URL           string
	GO_ENV           string
	DEBUG            bool
	MIGRATE_ON_START bool
	API_URL          string
	CLIENT_URL       string
	COOKIE_DOMAIN    string
	COOKIE_AGE_HOURS int
}

func APP() conf {
	return conf{
		PORT:      uint(getEnvInt("PORT", 8080)),
		JWT_TOKEN: getEnv("JWT_TOKEN", ""),
		DB_URL:    getEnv("DB_URL", ""),
		GO_ENV:    getEnv("GO_ENV", "development"),
		DEBUG:     getEnvBool("DEBUG", false),

		MIGRATE_ON_START: getEnvBool("MIGRATE_ON_START", false),
		API_URL:          getEnv("API_URL", ""),
		CLIENT_URL:       getEnv("CLIENT_URL", ""),
		COOKIE_DOMAIN:    getEnv("COOKIE_DOMAIN", ""),
		COOKIE_AGE_HOURS: getEnvInt("COOKIE_AGE_HOURS", 24),
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value

}

func getEnvInt(key string, defaultValue int) int {
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

func getEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
