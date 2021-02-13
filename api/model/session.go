package model

import "time"

// Session is ...
type Session struct {
	ID         string      `bson:"_id"`
	User       interface{} `bson:"user"`
	LastActive time.Time   `bson:"lastActive"`
}
