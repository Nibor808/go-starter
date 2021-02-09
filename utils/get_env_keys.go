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
	MailPass   string
	MailHost   string
}

// GetKeys returns the env keys
func GetKeys() (Keys, error) {
	devURL, exists := os.LookupEnv("DEV_URL")
	if !exists {
		return Keys{}, fmt.Errorf("cannot get DEV_URL from .env")
	}

	adminEmail, aEExists := os.LookupEnv("ADMIN_EMAIL")
	if !aEExists {
		return Keys{}, fmt.Errorf("cannot get ADMIN_EMAIL from .env")
	}

	mailPass, mPExists := os.LookupEnv("MAIL_PASS")
	if !mPExists {
		return Keys{}, fmt.Errorf("cannot get MAIL_PASS from .env")
	}

	mailHost, mPExists := os.LookupEnv("MAIL_HOST")
	if !mPExists {
		return Keys{}, fmt.Errorf("cannot get MAIL_HOST from .env")
	}

	return Keys{
		DevURL:     devURL,
		AdminEmail: adminEmail,
		MailPass:   mailPass,
		MailHost:   mailHost,
	}, nil
}
