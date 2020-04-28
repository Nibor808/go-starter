package controller

import (
	"context"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"go-starter/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func CreateSession(w http.ResponseWriter, userId primitive.ObjectID, ac AuthController) *http.Cookie {
	sID, _ := uuid.NewV4()

	sess := model.Session{
		Id:         sID.String(),
		User:       userId,
		LastActive: time.Now(),
	}

	if _, err := ac.db.Collection("sessions").InsertOne(context.TODO(), sess); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return &http.Cookie{
			Name:     "",
			Value:    "",
			MaxAge:   0,
			HttpOnly: false,
		}
	}

	return &http.Cookie{
		Name:     "go-starter",
		Value:    sID.String(),
		MaxAge:   600,
		HttpOnly: false,
	}
}

/* middleware */
func CheckSession(h httprouter.Handle, db *mongo.Database) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := model.Session{}

		c, err := r.Cookie("go-starter")
		if err != nil {
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
