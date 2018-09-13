package app

import "os"

var (
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_NAME     = os.Getenv("DB_NAME")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_SSL_MODE = os.Getenv("DB_SSL_MODE")

	EMAIL_NAME = os.Getenv("EMAIL_NAME")

	TEMP_EMAIL          = os.Getenv("TEMP_EMAIL")
	TEMP_EMAIL_PASSWORD = os.Getenv("TEMP_EMAIL_PASSWORD")
)
