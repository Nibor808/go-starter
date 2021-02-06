package utils_test

import (
	"go-starter/utils"
	"testing"
)

func TestCheckValidEmail(t *testing.T) {
	pass := "test@email.ca"
	isValid := utils.CheckValidEmail(pass)

	if !isValid {
		t.Errorf("Valid email did not pass. Expected: %t - got %t", true, isValid)
	}

	fails := []string{"testemail.ca", "test@emailca", "test@email", "test@email.commm"}

	for _, email := range fails {
		valid := utils.CheckValidEmail(email)

		if valid {
			t.Errorf("notValid email passed. Expected: %t - got %t", false, valid)
		}
	}
}
