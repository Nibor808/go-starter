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

type DataController struct {
	db *mongo.Database
}

func NewDataController(db *mongo.Database) *DataController {
	return &DataController{db}
}

func (fc DataController) SaveData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json charset=utf8")

	var data model.Data

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formResult, err := fc.db.Collection("forms").InsertOne(context.TODO(), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if fId, ok := formResult.InsertedID.(primitive.ObjectID); ok {
		data.Id = fId

		if err = json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (fc DataController) UpdateData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json charset=utf8")

	var data model.Data

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := fc.db.Collection("forms").UpdateOne(context.TODO(),
		bson.M{"_id": data.Id},
		bson.M{"$set": bson.M{"values": data.Values}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount != 0 {
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Data not updated", http.StatusInternalServerError)
	}
}
