package utils_test

import (
	"go-starter/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	date := utils.ParseDate("2006-01-02T15:04:05.000Z")
	assert.Equal(t, "Monday January 2 2006", date)
}

func TestCheckValidEmail(t *testing.T) {
	pass := "test@email.ca"
	isValid := utils.CheckValidEmail(pass)

	assert.Equal(t, isValid, true)

	fails := []string{"testemail.ca", "test@emailca", "test@email", "test@email.commm"}

	for _, email := range fails {
		isValid := utils.CheckValidEmail(email)

		assert.Equal(t, isValid, false)
	}
}
func TestEmailMessageConfig(t *testing.T) {
	mailArgs := utils.MailArgs{
		AdminEmail: "test@admin.ca",
		APIKey:     "TESTAPIKEY",
		Subject:    "Email Test",
		To:         "new_user@email.ca",
		HTML:       `<div>Testing email</div>`,
	}

	_, err := utils.MessageConfig(mailArgs)

	assert.Equal(t, err, nil)
}
