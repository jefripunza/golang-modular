package env

import "os"

func GetPort() string {
	value := os.Getenv("PORT")
	if value == "" {
		value = "3003"
	}
	return value
}

func GetCdnHostUrl() string {
	value := os.Getenv("CDN_HOST_URL")
	if value == "" {
		value = "https://cdn.cloufina.com"
	}
	return value
}
