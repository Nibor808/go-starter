package controller

import (
	"context"
	"go-starter/model"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateSession creates a Session
// saves it to the database
// returns pointer to a cookie
func CreateSession(w http.ResponseWriter, userID primitive.ObjectID, ac AuthController) *http.Cookie {
	sID := createsID()
	sess := createSession(sID, userID)
	saved := saveSessionToDB(w, ac, sess)
	c := createCookie(saved, sID)

	return c
}

func createsID() uuid.UUID {
	return uuid.NewV4()
}

func createSession(sID uuid.UUID, userID primitive.ObjectID) model.Session {
	sess := model.Session{
		ID:         sID.String(),
		User:       userID,
		LastActive: time.Now(),
	}

	return sess
}

func saveSessionToDB(w http.ResponseWriter, ac AuthController, sess model.Session) bool {
	if _, err := ac.db.Collection("sessions").InsertOne(context.TODO(), sess); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}

func createCookie(saved bool, sID uuid.UUID) *http.Cookie {
	if saved {
		return &http.Cookie{
			Name:     "go-starter",
			Value:    sID.String(),
			MaxAge:   600,
			HttpOnly: false,
		}
	}

	return &http.Cookie{
		Name:     "",
		Value:    "",
		MaxAge:   0,
		HttpOnly: false,
	}
}
