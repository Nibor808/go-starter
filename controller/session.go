package controller

import (
	"context"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"go-starter/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func CreateSession(w http.ResponseWriter, userId primitive.ObjectID, uc UserController) *http.Cookie {
	sID, _ := uuid.NewV4()

	s := model.Session{
		Id:         sID.String(),
		User:       userId,
		LastActive: time.Now(),
	}

	_, err := uc.db.Collection("sessions").InsertOne(context.TODO(), s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Session not created")
	}

	return &http.Cookie{
		Name:     "go-starter",
		Value:    sID.String(),
		MaxAge:   600,
		HttpOnly: false,
	}
}

func CheckSession(h httprouter.Handle, db *mongo.Database) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c, err := r.Cookie("go-starter")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		} else {
			_, err = db.Collection("sessions").Find(context.TODO(), bson.M{"_id": c.Value})
			if err != nil {
				http.Error(w, "you must be signed in to view this info", http.StatusUnauthorized)
			}

			h(w, r, p)
		}
	}
}
