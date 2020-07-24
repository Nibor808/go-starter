package controller

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go-starter/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type FormController struct {
	db *mongo.Database
}

func NewFormController(db *mongo.Database) *FormController {
	return &FormController{db}
}

func (fc FormController) SaveForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json charset=utf8")

	var f model.Form

	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formResult, err := fc.db.Collection("forms").InsertOne(context.TODO(), f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if fId, ok := formResult.InsertedID.(primitive.ObjectID); ok {
		f.Id = fId

		if err = json.NewEncoder(w).Encode(f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (fc FormController) UpdateForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json charset=utf8")

	var f model.Form

	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := fc.db.Collection("forms").UpdateOne(context.TODO(),
		bson.M{"_id": f.Id},
		bson.M{"$set": bson.M{"values": f.Values}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount != 0 {
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Form not updated", http.StatusInternalServerError)
	}
}
