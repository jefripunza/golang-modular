package config

import (
	"core/env"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  env.GetGoogleRedirectURL(),
		ClientID:     env.GetClientID(),
		ClientSecret: env.GetClientSecret(),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	GoogleOauthStateString = uuid.New().String() // You should generate a random state string for CSRF protection
)
