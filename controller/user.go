package controller

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go-starter/model"
	"go-starter/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserController struct {
	db *mongo.Database
}

func NewUserController(db *mongo.Database) *UserController {
	return &UserController{db}
}

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
			log.Fatal("User cursor.Close error:", cerr)
		}
	}()

	cursorErr := cursor.All(context.TODO(), &results)
	if cursorErr != nil {
		http.Error(w, cursorErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
	}
}

func (uc UserController) User(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	result := model.User{}

	docId, err := primitive.ObjectIDFromHex(p.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = uc.db.Collection("users").FindOne(context.TODO(), bson.M{"_id": docId}).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
	}
}

func (uc UserController) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	p := r.FormValue("password")

	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
	}

	u := model.User{
		Email: r.FormValue("email"),
		Password: string(bs),
	}

	_, err = uc.db.Collection("users").InsertOne(context.TODO(), u)
	if err != nil {
		var merr mongo.WriteException
		merr = err.(mongo.WriteException)
		errCode := merr.WriteErrors[0].Code

		if errCode == 11000 {
			data := map[string]string{"error": "That email is already in use."}
			err = json.NewEncoder(w).Encode(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		mailSent := utils.SendMail("Test", r.FormValue("email"), "Test content")
		var data map[string]string

		if mailSent {
			data = map[string]string{"ok": "Email sent"}
		} else {
			data = map[string]string{"error": "Unable to send email. Please contact support."}
		}

		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
