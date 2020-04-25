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

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uc UserController) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var returnData map[string]string

	pass := r.FormValue("password")

	bs, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := model.User{
		Email:    r.FormValue("email"),
		Password: string(bs),
	}

	result, err := uc.db.Collection("users").InsertOne(context.TODO(), u)
	if err != nil {
		var merr mongo.WriteException

		merr = err.(mongo.WriteException)
		errCode := merr.WriteErrors[0].Code

		if errCode == 11000 {
			w.WriteHeader(http.StatusForbidden)
			returnData = map[string]string{"error": "That email is already in use."}

			err = json.NewEncoder(w).Encode(returnData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		mailSent := utils.SendMail("Go Starter", r.FormValue("email"), "Welcome To Go Starter")

		if mailSent {
			if oId, ok := result.InsertedID.(primitive.ObjectID); ok {
				c := CreateSession(w, oId, uc)
				http.SetCookie(w, c)

				w.WriteHeader(http.StatusCreated)
				returnData = map[string]string{"ok": "Email sent"}
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			returnData = map[string]string{"error": "Unable to send email. Please contact support."}
		}

		err = json.NewEncoder(w).Encode(returnData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (uc UserController) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	email := r.FormValue("email")
	pass := r.FormValue("password")
	user := model.User{}

	err := uc.db.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "username and/or password incorrect", http.StatusUnauthorized)
		return
	} else {
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
		if passErr != nil {
			http.Error(w, "username and/or password incorrect", http.StatusUnauthorized)
			return
		}

		userId, objectFromHexErr := primitive.ObjectIDFromHex(user.Id)
		if objectFromHexErr != nil {
			http.Error(w, objectFromHexErr.Error(), http.StatusInternalServerError)
			return
		}

		c := CreateSession(w, userId, uc)
		http.SetCookie(w, c)

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (uc UserController) SignOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result bson.M

	c, err := r.Cookie("go-starter")
	if err != nil {
		http.Error(w, "cannot get cookie", http.StatusInternalServerError)
		return
	}

	err = uc.db.Collection("sessions").FindOneAndDelete(context.TODO(), bson.M{"_id": c.Value}).Decode(&result)
	if err != nil {
		http.Error(w, "session not deleted", http.StatusInternalServerError)
		return
	}

	c = &http.Cookie{
		Name:   "go-starter",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, c)
}
