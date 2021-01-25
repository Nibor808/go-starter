package controller

import (
	"context"
	"fmt"
	"go-starter/model"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey []byte

// init gets the JWT_KEY from .env
func init() {
	key, exists := os.LookupEnv("JWT_KEY")
	if !exists {
		log.Println("Cannot get JWT_KEY from .env")
	} else {
		jwtKey = []byte(key)
	}
}

// myClaims is...
type myClaims struct {
	jwt.StandardClaims
}

// Valid validates the jwt signature
func (c myClaims) Valid() error {
	if !c.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token has expired")
	}

	if c.StandardClaims.Id == "" {
		return fmt.Errorf("token invalid")
	}

	return nil
}

// GetCookie creates a Session
// saves it to the database
// returns pointer to a cookie
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

func ParseToken(token string) (*myClaims, error) {
	tokenAfter, err := jwt.ParseWithClaims(token, &myClaims{}, func(tBefore *jwt.Token) (interface{}, error) {
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

	return tokenAfter.Claims.(*myClaims), nil
}

func createJwt(sID string) (string, error) {
	mc := myClaims{
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

func createUID() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error in createSID getting uuid: %w", err)
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

func createCookie(saved bool, ss string) *http.Cookie {
	if saved {
		return &http.Cookie{
			Name:     "go-starter",
			Value:    ss,
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
