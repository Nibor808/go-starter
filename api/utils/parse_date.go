package utils

import (
	"fmt"
	"log"
	"time"
)

// ParseDate parses a date string
// eg "2006-01-02T15:04:05.000Z"
// into "Monday January 2 2006"
func ParseDate(d string) string {
	layout := "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, d)
	if err != nil {
		log.Fatalln("Date parse failed", err)
	}

	return fmt.Sprintf("%s %s %d %d", t.Weekday(), t.Month(), t.Day(), t.Year())
}
