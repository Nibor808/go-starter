package model

import (
	"time"
)

type Session struct {
	Id         string      `bson:"_id"`
	User       interface{} `bson:"user"`
	LastActive time.Time   `bson:"lastActive"`
}
