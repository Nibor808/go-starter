package utils

import "testing"

func TestCheckValidEmail(t *testing.T) {
	validEmail := "test@gmail.com"
	invalidEmails := []string{"test@gmail.co", "test@gmailcom", "testgmail.com", "@gmail.com"}

	isValid := CheckValidEmail(validEmail)

	if !isValid {
		t.Errorf("Email %s is invalid", validEmail)
	}

	for _, email := range invalidEmails {
		isValid = CheckValidEmail(email)

		if isValid {
			t.Errorf("Email %s is invalid but returned as valid", email)
		}
	}
}
