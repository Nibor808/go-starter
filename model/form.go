package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Form struct {
	Id     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Values interface{}        `json:"values" bson:"values"`
}
