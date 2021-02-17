package model

// User is ...
type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	IsActive bool   `json:"isActive" bson:"isActive"`
	IsAdmin  bool   `json:"isAdmin" bson:"isAdmin"`
}
