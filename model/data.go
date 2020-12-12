package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Data is ...
type Data struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Values   interface{}        `json:"values" bson:"values"`
	FileName string
}
