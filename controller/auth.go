package controller

import (
	"context"
	"encoding/json"
	"go-starter/model"
	"go-starter/utils"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// AuthController is ...
type AuthController struct {
	db *mongo.Database
}

// NewAuthController is ...
func NewAuthController(db *mongo.Database) *AuthController {
	return &AuthController{db}
}

// SignUpEmail checks for email validity
// checks that email doesn't already exist in the database
// creates a user with the email given
// sends a verification email
// sets a cookie
func (ac AuthController) SignUpEmail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	email := r.FormValue("email")

	u := model.User{
		Email: email,
	}

	isValid := utils.CheckValidEmail(email)
	if !isValid {
		http.Error(w, "Invalid Email", http.StatusInternalServerError)
		return
	}

	// go ahead and insert the new user
	userResult, err := ac.db.Collection("users").InsertOne(context.TODO(), u)
	if err != nil {
		var merr = err.(mongo.WriteException)
		errCode := merr.WriteErrors[0].Code

		if errCode == 11000 {
			http.Error(w, "That email is already in use", http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := model.Token{
		UserID:       userResult.InsertedID,
		CreationTime: time.Now(),
	}

	tokenResult, err := ac.db.Collection("tokens").InsertOne(context.TODO(), t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	keys, err := utils.GetKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var htmlText strings.Builder

	if tID, ok := tokenResult.InsertedID.(primitive.ObjectID); ok {
		if uID, ok := userResult.InsertedID.(primitive.ObjectID); ok {
			htmlText.WriteString("Welcome To Go Starter! Follow " +
				"the <a href=" + keys.DevURL + "/confirmemail/" +
				tID.Hex() + "/" + uID.Hex() + ">Link</a> or paste " +
				"this into your browser's address bar: " + keys.DevURL +
				"/confirmemail/" + tID.Hex() + "/" + uID.Hex())
		}
	}

	mailArgs := utils.MailArgs{
		AdminEmail: keys.AdminEmail,
		APIKey:     keys.APIKey,
		Subject:    "Go Starter",
		To:         email,
		HTML:       htmlText.String(),
	}

	err = utils.SendMail(mailArgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if uID, ok := userResult.InsertedID.(primitive.ObjectID); ok {
		c := ac.GetCookie(w, uID)
		http.SetCookie(w, c)

		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte("Email sent to " + email))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ConfirmEmail confirms the verification email
// deletes the Token
// deletes the user if the token has expired
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
	if err != nil { // no token
		http.Error(w, "Token expired.", http.StatusUnauthorized)

		err := ac.db.Collection("users").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
		if err != nil { // no user
			http.Error(w, "Link expired. Sign up at /signupemail, { email: \"your-email\" }.", http.StatusUnauthorized)
			return
		} else if len(user.Password) == 0 { // no token / have user / no password
			http.Error(w, "Link expired. Sign up at /signupemail, { email: \"your-email\" }.", http.StatusUnauthorized)

			// delete user so sign up can proceed again without unique email conflict
			err = ac.db.Collection("users").FindOneAndDelete(context.TODO(), bson.M{"_id": userID}).Decode(&deletedDoc)
			if err != nil {
				http.Error(w, "Unable to delete user", http.StatusUnauthorized)
			}

			return
		} else {
			// no token / have user / have password
			http.Error(w, "User is already signed up. Go ahead and sign in.", http.StatusUnauthorized)
			return
		}
	} else {
		err := ac.db.Collection("users").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Email verified. Go ahead and save a password at /signuppassword, { password: \"your-password\" }"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// SignUpPassword collects a password from a form
// checks that the cookie is valid
// updates the user by adding the hashed password
func (ac AuthController) SignUpPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	pass := r.FormValue("password")
	sess := model.Session{}

	c, err := r.Cookie("go-starter")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mc, err := ParseToken(c.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ac.db.Collection("sessions").FindOne(context.TODO(), bson.M{"_id": mc.StandardClaims.Id}).Decode(&sess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := ac.db.Collection("users").UpdateOne(context.TODO(),
		bson.M{"_id": sess.User},
		bson.M{"$set": bson.M{"password": string(bs), "isActive": true}})
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

// SignIn signs in
// sets a cookie
// returns the user
func (ac AuthController) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	email := r.FormValue("email")
	pass := r.FormValue("password")
	user := model.User{}

	err := ac.db.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "Username and/or password incorrect", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		http.Error(w, "Username and/or password incorrect", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := ac.GetCookie(w, userID)
	http.SetCookie(w, c)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SignOut signs out the user
// checks that the cookie is valid
// deletes the session
// resets the cookie
func (ac AuthController) SignOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result bson.M

	c, err := r.Cookie("go-starter")
	if err != nil {
		http.Error(w, "Cannot get cookie", http.StatusInternalServerError)
		return
	}
	mc, err := ParseToken(c.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = ac.db.Collection("sessions").FindOneAndDelete(context.TODO(), bson.M{"_id": mc.StandardClaims.Id}).Decode(&result)
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
