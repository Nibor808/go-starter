package controller

import (
	"api/model"
	"api/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey []byte

// init gets the JWT_KEY from .env
func init() {
	keys, err := utils.GetKeys()
	if err != nil {
		log.Println(err)
	}

	jwtKey = []byte(keys.JWTKey)
}

// MyClaims is...
// using ExpiresAt and Id in jwt.StandardClaims
type MyClaims struct {
	jwt.StandardClaims
}

// Valid validates the JWT token
func (c MyClaims) Valid() error {
	if !c.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token has expired")
	}

	if c.StandardClaims.Id == "" {
		return fmt.Errorf("token invalid")
	}

	return nil
}

// GetCookie returns the created cookie
func (ac AuthController) GetCookie(w http.ResponseWriter, userID primitive.ObjectID) *http.Cookie {
	sID, err := createUID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	ss, err := createJwt(sID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	sess := createSession(sID, userID)
	saved := saveSessionToDB(w, ac, sess)
	c := createCookie(saved, ss)

	return c
}

// ParseToken takes in a JWT token from a cookie
// checks that the signing method is the same as
// the one used to create the token
// calls the Valid function on the verified token
// asserts that the token.Claims are of type *MyClaims
// returns the token and an error
func ParseToken(token string) (*MyClaims, error) {
	tokenAfter, err := jwt.ParseWithClaims(token, &MyClaims{}, func(tBefore *jwt.Token) (interface{}, error) {
		if tBefore.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("token invalid")
		}

		return jwtKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if !tokenAfter.Valid {
		return nil, fmt.Errorf("token invalid: %w", err)
	}

	return tokenAfter.Claims.(*MyClaims), nil
}

// createJWT takes in a sID
// creates a JWT token with the sID as jwt.StandardClaims.Id
// returns a signed string and an error
func createJwt(sID string) (string, error) {
	mc := MyClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Id:        sID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, mc)
	ss, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error getting SignedString: %w", err)
	}

	return ss, nil
}

// createUID returns a new v4 uuid and an error
func createUID() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error in createSID getting uuid: %w", err)
	}

	return uid.String(), nil
}

// createSession takes in a sID and a user id
// adds those values to the Session
// sets LastActive to time.Now()
// returns the Session
func createSession(sID string, userID primitive.ObjectID) model.Session {
	sess := model.Session{
		ID:         sID,
		User:       userID,
		LastActive: time.Now(),
	}

	return sess
}

// saveSessionToDB saves the Session to the database
func saveSessionToDB(w http.ResponseWriter, ac AuthController, sess model.Session) bool {
	if _, err := ac.db.Collection("sessions").InsertOne(context.TODO(), sess); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}

// createCookie returns a pointer to a cookie
// with the JWT token as the value if the Session
// was saved to the database
// otherwise returns a pointer to a cookie with no values
func createCookie(saved bool, ss string) *http.Cookie {
	if saved {
		return &http.Cookie{
			Name:     "go-starter",
			Value:    ss,
			MaxAge:   600,
			HttpOnly: true,
		}
	}

	return &http.Cookie{
		Name:     "",
		Value:    "",
		MaxAge:   0,
		HttpOnly: true,
	}
}
