package utils

import (
	"fmt"
	"os"
)

// Keys is ...
type Keys struct {
	DevURL     string
	AdminEmail string
	APIKey     string
}

// GetKeys returns the env keys
func GetKeys() (Keys, error) {
	devURL, exists := os.LookupEnv("DEV_URL")
	if !exists {
		return Keys{}, fmt.Errorf("Cannot get DEV_URL from .env")
	}

	adminEmail, sEExists := os.LookupEnv("ADMIN_EMAIL")
	if !sEExists {
		return Keys{}, fmt.Errorf("Cannot get ADMIN_EMAIL from .env")
	}

	apiKey, kExists := os.LookupEnv("SENDGRID_API_KEY")
	if !kExists {
		return Keys{}, fmt.Errorf("Cannot get SENDGRID_API_KEY from .env")
	}

	return Keys{
		DevURL:     devURL,
		AdminEmail: adminEmail,
		APIKey:     apiKey,
	}, nil
}
