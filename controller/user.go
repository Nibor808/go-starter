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
	"strings"
	"time"
)

type UserController struct {
	db *mongo.Database
}

func NewUserController(db *mongo.Database) *UserController {
	return &UserController{db}
}

/* return all users */
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

	if cursorErr := cursor.All(context.TODO(), &results); cursorErr != nil {
		http.Error(w, cursorErr.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/* returns user from current session */
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

/* create user with email from a form and send verification email */
func (uc UserController) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	user := model.User{
		Email: r.FormValue("email"),
	}

	userResult, err := uc.db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		var merr mongo.WriteException

		merr = err.(mongo.WriteException)
		errCode := merr.WriteErrors[0].Code

		if errCode == 11000 {
			http.Error(w, "That email is already in use", http.StatusConflict)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var htmlText strings.Builder

		t := model.Token{
			UserId:       userResult.InsertedID,
			CreationTime: time.Now(),
		}

		tokenResult, err := uc.db.Collection("tokens").InsertOne(context.TODO(), t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if tId, ok := tokenResult.InsertedID.(primitive.ObjectID); ok {
			if uId, ok := userResult.InsertedID.(primitive.ObjectID); ok {
				htmlText.WriteString("Welcome To Go Starter! Follow " +
					"the <a href=http://localhost:8080/verifyemail/" + tId.Hex() + "/" + uId.Hex() + ">Link</a> or paste " +
					"this into your browser's address bar: localhost:8080/" + tId.Hex() + "/" + uId.Hex())
			}
		}

		mailSent := utils.SendMail("Go Starter", r.FormValue("email"), htmlText.String())

		if mailSent {
			if uId, ok := userResult.InsertedID.(primitive.ObjectID); ok {
				c := CreateSession(w, uId, uc)
				http.SetCookie(w, c)

				w.WriteHeader(http.StatusCreated)
				_, err := w.Write([]byte("Email sent"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		} else {
			http.Error(w, "Email not sent", http.StatusInternalServerError)
		}
	}
}

/* confirm verification email and delete Token, delete user if token has expired */
func (uc UserController) ConfirmVerificationEmail(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	user := model.User{}
	var deletedDoc bson.M

	token, err := primitive.ObjectIDFromHex(p.ByName("token"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := primitive.ObjectIDFromHex(p.ByName("userId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = uc.db.Collection("tokens").FindOneAndDelete(context.TODO(), bson.M{"_id": token}).Decode(&deletedDoc)
	if err != nil {
		http.Error(w, "Token expired. Sign up again.", http.StatusUnauthorized)

		err = uc.db.Collection("users").FindOneAndDelete(context.TODO(), bson.M{"_id": userId}).Decode(&deletedDoc)
		if err != nil {
			http.Error(w, "Unable to delete user", http.StatusUnauthorized)
			return
		}
	} else {
		err := uc.db.Collection("users").FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Email verified"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/* collect password from a form and update user */
func (uc UserController) CollectPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	pass := r.FormValue("password")
	sess := model.Session{}
	c, err := r.Cookie("go-starter")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = uc.db.Collection("sessions").FindOne(context.TODO(), bson.M{"_id": c.Value}).Decode(&sess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := uc.db.Collection("users").UpdateOne(context.TODO(),
		bson.M{"_id": sess.User},
		bson.M{"$set": bson.M{"password": string(bs)}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount != 0 {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Password confirmed. User updated."))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "User not updated", http.StatusInternalServerError)
	}
}

/* sign in and return user */
func (uc UserController) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	email := r.FormValue("email")
	pass := r.FormValue("password")
	user := model.User{}

	err := uc.db.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "username and/or password incorrect", http.StatusUnauthorized)
		return
	} else {
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
			http.Error(w, "Username and/or password incorrect", http.StatusUnauthorized)
			return
		}

		userId, err := primitive.ObjectIDFromHex(user.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c := CreateSession(w, userId, uc)
		http.SetCookie(w, c)

		if err = json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/* sign out user, delete session and reset cookie */
func (uc UserController) SignOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result bson.M

	c, err := r.Cookie("go-starter")
	if err != nil {
		http.Error(w, "Cannot get cookie", http.StatusInternalServerError)
		return
	}

	err = uc.db.Collection("sessions").FindOneAndDelete(context.TODO(), bson.M{"_id": c.Value}).Decode(&result)
	if err != nil {
		http.Error(w, "Session not deleted", http.StatusInternalServerError)
		return
	}

	c = &http.Cookie{
		Name:   "go-starter",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, c)
}
