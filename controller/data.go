package controller

import (
	"context"
	"encoding/json"
	"go-starter/model"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataController struct {
	db *mongo.Database
}

func NewDataController(db *mongo.Database) *DataController {
	return &DataController{db}
}

func (dc DataController) AllData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json charset=utf8")

	var data []model.Data

	cursor, findError := dc.db.Collection("data").Find(context.TODO(), bson.D{{}})
	if findError != nil {
		http.Error(w, findError.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if cerr := cursor.Close(context.TODO()); cerr != nil {
			log.Fatal("Data cursor. Close error:", cerr)
		}
	}()

	if cursorErr := cursor.All(context.TODO(), &data); cursorErr != nil {
		http.Error(w, cursorErr.Error(), http.StatusInternalServerError)
		return
	}

	if len(data) <= 0 {
		http.Error(w, "No data available.", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (dc DataController) SaveData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json charset=utf8")

	var data model.Data

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formResult, err := dc.db.Collection("data").InsertOne(context.TODO(), data)
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

func (dc DataController) UpdateData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json charset=utf8")

	var data model.Data

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := dc.db.Collection("data").UpdateOne(context.TODO(),
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
