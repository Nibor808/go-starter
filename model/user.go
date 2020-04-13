package model

type User struct {
	Id            string `json:"id,omitempty" bson:"_id,omitempty"`
	Email         string `json:"email" bson:"email"`
	Password      string `json:"password" bson:"password"`
	IsActive      bool   `json:"isActive" bson:"isActive"`
	IsAdmin       bool   `json:"isAdmin" bson:"isAdmin"`
	StripeAccount string `json:"stripe_account" bson:"stripe_account"`
	Plan          string `json:"plan" bson:"plan"`
	Subscription  string `json:"subscription" bson:"subscription"`
}
