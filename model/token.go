package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token is ...
type Token struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       interface{}        `bson:"UserID"`
	CreationTime time.Time          `bson:"creationTime"`
}
