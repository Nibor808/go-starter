package middleware

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-starter/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func CheckSession(h httprouter.Handle, db *mongo.Database) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := model.Session{}

		c, err := r.Cookie("go-starter")
		if err != nil {
			fmt.Println("HERE-------SESS")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		} else {
			err = db.Collection("sessions").FindOneAndUpdate(
				context.TODO(),
				bson.M{"_id": c.Value},
				bson.M{"$set": bson.M{"lastActive": time.Now()}}).Decode(&sess)
			if err != nil {
				http.Error(w, "Session expired", http.StatusUnauthorized)
				return
			}

			h(w, r, p)
		}
	}
}
