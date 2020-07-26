package controller

import (
	"context"
	"encoding/json"
	"go-starter/model"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserController is ...
type UserController struct {
	db *mongo.Database
}

// NewUserController is ...
func NewUserController(db *mongo.Database) *UserController {
	return &UserController{db}
}

// AllUsers returns all users in db
func (uc UserController) AllUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var results []model.User

	cursor, findError := uc.db.Collection("users").Find(context.TODO(), bson.D{{}})
	if findError != nil {
		http.Error(w, findError.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if cerr := cursor.Close(context.TODO()); cerr != nil {
			log.Fatal("User cursor. Close error:", cerr)
		}
	}()

	if cursorErr := cursor.All(context.TODO(), &results); cursorErr != nil {
		http.Error(w, cursorErr.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// User returns the user info from the current session
func (uc UserController) User(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	user := model.User{}
	sess := model.Session{}

	c, _ := r.Cookie("go-starter")

	err := uc.db.Collection("sessions").FindOne(context.TODO(), bson.M{"_id": c.Value}).Decode(&sess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = uc.db.Collection("users").FindOne(context.TODO(), bson.M{"_id": sess.User}).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
