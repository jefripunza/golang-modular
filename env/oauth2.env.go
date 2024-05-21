package env

import "os"

func GetClientID() string {
	value := os.Getenv("SECRET_KEY")
	if value == "" {
		value = "your_google_client_id"
	}
	return value
}

func GetClientSecret() string {
	value := os.Getenv("SECRET_KEY")
	if value == "" {
		value = "your_google_client_secret"
	}
	return value
}

func GetGoogleRedirectURL() string {
	value := os.Getenv("SECRET_KEY")
	if value == "" {
		value = "http://localhost:3003/auth/google/callback"
	}
	return value
}
