package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

// Token is ...
type Token struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       interface{}        `bson:"UserID"`
	CreationTime time.Time          `bson:"creationTime"`
}
