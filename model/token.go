package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Token struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	UserId       interface{}        `bson:"userId"`
	CreationTime time.Time          `bson:"creationTime"`
}
