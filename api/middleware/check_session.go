package middleware

import (
	"api/controller"
	"api/model"
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CheckSession checks that that cookie exists
// validates the jwt token
// returns the handler function given
// updates the cookie lastActive if it exists
func CheckSession(h httprouter.Handle, db *mongo.Database) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := model.Session{}

		c, err := r.Cookie("go-starter")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		mc, err := controller.ParseToken(c.Value)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		err = db.Collection("sessions").FindOneAndUpdate(
			context.TODO(),
			bson.M{"_id": mc.StandardClaims.Id},
			bson.M{"$set": bson.M{"lastActive": time.Now()}}).Decode(&sess)
		if err != nil {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		h(w, r, p)
	}
}
