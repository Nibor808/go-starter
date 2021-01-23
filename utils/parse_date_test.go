package utils

import "testing"

func TestParseDate(t *testing.T) {
	date := ParseDate("2006-01-02T15:04:05.000Z")

	if date != "Monday January 2 2006" {
		t.Errorf("Date was incorrect, got: %s, want: %s", date, "Monday January 2 2006")
	}
}