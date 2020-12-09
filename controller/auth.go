package controller

import (
	"context"
	"encoding/json"
	"go-starter/model"
	"go-starter/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db *mongo.Database
}

func NewAuthController(db *mongo.Database) *AuthController {
	return &AuthController{db}
}

/* create user with email from a form and send verification email */
func (ac AuthController) SignUpEmail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	u := model.User{
		Email: r.FormValue("email"),
	}

	userResult, err := ac.db.Collection("users").InsertOne(context.TODO(), u)
	if err != nil {
		var merr = err.(mongo.WriteException)
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
			UserID:       userResult.InsertedID,
			CreationTime: time.Now(),
		}

		tokenResult, err := ac.db.Collection("tokens").InsertOne(context.TODO(), t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dev_url, exists := os.LookupEnv("DEV_URL")
		if !exists {
			http.Error(w, "Cannot find DEV_URL", http.StatusInternalServerError)
			return
		}

		if tId, ok := tokenResult.InsertedID.(primitive.ObjectID); ok {
			if uId, ok := userResult.InsertedID.(primitive.ObjectID); ok {
				htmlText.WriteString("Welcome To Go Starter! Follow " +
					"the <a href=" + dev_url + "/confirmemail/" +
					tId.Hex() + "/" + uId.Hex() + ">Link</a> or paste " +
					"this into your browser's address bar: " + dev_url +
					"/confirmemail/" + tId.Hex() + "/" + uId.Hex())
			}
		}

		mailSent := utils.SendMail("Go Starter", r.FormValue("email"), htmlText.String())

		if mailSent {
			if uId, ok := userResult.InsertedID.(primitive.ObjectID); ok {
				c := CreateSession(w, uId, ac)
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
func (ac AuthController) ConfirmEmail(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	user := model.User{}
	var deletedDoc bson.M

	token, err := primitive.ObjectIDFromHex(p.ByName("token"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, err := primitive.ObjectIDFromHex(p.ByName("userID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ac.db.Collection("tokens").FindOneAndDelete(context.TODO(), bson.M{"_id": token}).Decode(&deletedDoc)
	if err != nil {
		http.Error(w, "Token expired. Sign up again.", http.StatusUnauthorized)

		err = ac.db.Collection("users").FindOneAndDelete(context.TODO(), bson.M{"_id": userID}).Decode(&deletedDoc)
		if err != nil {
			http.Error(w, "Unable to delete user", http.StatusUnauthorized)
			return
		}
	} else {
		err := ac.db.Collection("users").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
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
func (ac AuthController) SignUpPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	err = ac.db.Collection("sessions").FindOne(context.TODO(), bson.M{"_id": c.Value}).Decode(&sess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := ac.db.Collection("users").UpdateOne(context.TODO(),
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
func (ac AuthController) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	email := r.FormValue("email")
	pass := r.FormValue("password")
	user := model.User{}

	err := ac.db.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "Username and/or password incorrect", http.StatusUnauthorized)
		return
	} else {
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
			http.Error(w, "Username and/or password incorrect", http.StatusUnauthorized)
			return
		}

		UserID, err := primitive.ObjectIDFromHex(user.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c := CreateSession(w, UserID, ac)
		http.SetCookie(w, c)

		if err = json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/* sign out user, delete session and reset cookie */
func (ac AuthController) SignOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result bson.M

	c, err := r.Cookie("go-starter")
	if err != nil {
		http.Error(w, "Cannot get cookie", http.StatusInternalServerError)
		return
	}

	err = ac.db.Collection("sessions").FindOneAndDelete(context.TODO(), bson.M{"_id": c.Value}).Decode(&result)
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

	_, err = w.Write([]byte("Signed out"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
