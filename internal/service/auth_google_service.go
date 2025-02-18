package service

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}

// Obtener URL de autenticación
func GetGoogleAuthURL() string {
	return GoogleOAuthConfig.AuthCodeURL("randomstate")
}

// Obtener token de usuario
func GetGoogleUserInfo(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	return GoogleOAuthConfig.Exchange(ctx, code)
}
