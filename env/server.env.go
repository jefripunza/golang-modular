package env

import (
	"log"
	"os"
	"strconv"
)

func GetServerName() string {
	value := os.Getenv("SERVER_NAME")
	if value == "" {
		value = "P34C3_KHYREIN"
	}
	return value
}

func GetServerPort() string {
	value := os.Getenv("SERVER_PORT")
	if value == "" {
		value = "3003"
	}
	return value
}

func GetSecretKey() string {
	value := os.Getenv("SECRET_KEY")
	if value == "" {
		value = "your_secret_key"
	}
	return value
}

func GetMaxLoginAttempts() int64 {
	value := os.Getenv("MAX_LOGIN_ATTEMPTS")
	if value == "" {
		return 3
	}
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("Error parsing MAX_LOGIN_ATTEMPTS: %s. Using default value 3.", err)
		return 3
	}
	return num
}

func GetCdnHostUrl() string {
	value := os.Getenv("CDN_HOST_URL")
	if value == "" {
		value = "https://cdn.cloufina.com"
	}
	return value
}
