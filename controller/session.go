package controller

import (
	"context"
	"fmt"
	"go-starter/model"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCookie creates a Session
// saves it to the database
// returns pointer to a cookie
func (ac AuthController) GetCookie(w http.ResponseWriter, userID primitive.ObjectID) *http.Cookie {
	sID, err := createSID()
	if err != nil {
		log.Println("error creating sID in GetCookie")
	}

	sess := createSession(sID, userID)
	saved := saveSessionToDB(w, ac, sess)
	c := createCookie(saved, sID)

	return c
}

func createSID() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error in genereateSid getting uuid: %w", err)
	}

	return uid.String(), nil
}

func createSession(sID string, userID primitive.ObjectID) model.Session {
	sess := model.Session{
		ID:         sID,
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

func createCookie(saved bool, sID string) *http.Cookie {
	if saved {
		return &http.Cookie{
			Name:     "go-starter",
			Value:    sID,
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
