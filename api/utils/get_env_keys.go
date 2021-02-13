package utils

import (
	"fmt"
	"os"
)

// Keys is ...
type Keys struct {
	DevURL     string
	DBConn     string
	DBName     string
	AdminEmail string
	MailPass   string
	MailHost   string
	JWTKey     string
}

// GetKeys returns the env keys
func GetKeys() (Keys, error) {
	devURL, exists := os.LookupEnv("DEV_URL")
	if !exists {
		return Keys{}, fmt.Errorf("cannot get DEV_URL from .env")
	}

	dbConn, exists := os.LookupEnv("DB_CONN")
	if !exists {
		return Keys{}, fmt.Errorf("cannot get DB_CONN from .env")
	}

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		return Keys{}, fmt.Errorf("cannot get DB_NAME from .env")
	}

	adminEmail, exists := os.LookupEnv("ADMIN_EMAIL")
	if !exists {
		return Keys{}, fmt.Errorf("cannot get ADMIN_EMAIL from .env")
	}

	mailPass, exists := os.LookupEnv("MAIL_PASS")
	if !exists {
		return Keys{}, fmt.Errorf("cannot get MAIL_PASS from .env")
	}

	mailHost, exists := os.LookupEnv("MAIL_HOST")
	if !exists {
		return Keys{}, fmt.Errorf("cannot get MAIL_HOST from .env")
	}

	jwtKey, exists := os.LookupEnv("JWT_KEY")
	if !exists {
		return Keys{}, fmt.Errorf("cannot get JWT_KEY from .env")
	}

	return Keys{
		DevURL:     devURL,
		DBConn:     dbConn,
		DBName:     dbName,
		AdminEmail: adminEmail,
		MailPass:   mailPass,
		MailHost:   mailHost,
		JWTKey:     jwtKey,
	}, nil
}
