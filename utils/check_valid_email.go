package utils

import (
	"net"
	"regexp"
	"strings"
)

func CheckValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(email) < 3 || len(email) > 254 {
		return false
	}

	if !emailRegex.MatchString(email) {
		return false
	}

	parts := strings.Split(email, "@")
	mx_records, err := net.LookupMX(parts[1])

	if err != nil || len(mx_records) == 0 {
		return false
	}

	return true
}
