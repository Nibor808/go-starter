package controller

import (
	"context"
	"go-starter/model"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSession(w http.ResponseWriter, userID primitive.ObjectID, ac AuthController) *http.Cookie {
	sID := uuid.NewV4()

	sess := model.Session{
		Id:         sID.String(),
		User:       userID,
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
