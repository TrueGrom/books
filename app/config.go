package app

import "os"

var defaultValues = map[string]string{
	"DB_HOST":     "127.0.0.1",
	"DB_PORT":     "5432",
	"DB_SSL_MODE": "disabled",
}

func getenvDefault(key string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValues[key]
	}
	return value
}

var (
	DB_HOST     = getenvDefault("DB_HOST")
	DB_PORT     = getenvDefault("DB_PORT")
	DB_NAME     = os.Getenv("DB_NAME")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_SSL_MODE = getenvDefault("DB_SSL_MODE")

	EMAIL_NAME = os.Getenv("EMAIL_NAME")

	TEMP_EMAIL          = os.Getenv("TEMP_EMAIL")
	TEMP_EMAIL_PASSWORD = os.Getenv("TEMP_EMAIL_PASSWORD")
)
