package utils

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

func GetMongoSession() *mongo.Database {
	dbConn, connExists := os.LookupEnv("DB_CONN")
	if !connExists {
		log.Fatal("Missing db connection string from .env")
	}

	dbName, nameExists := os.LookupEnv("DB_NAME")
	if !nameExists {
		log.Fatal("Missing db name from .env")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConn))
	if err != nil {
		log.Fatalln(err)
	}

	db := client.Database(dbName)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}

	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)

	emailIndex, err := db.Collection("users").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		}, opts)
	if err != nil {
		log.Fatalln(err)
	}

	ttl := int32(60)
	tokenIndex, err := db.Collection("tokens").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{"creationTime": 1},
			Options: &options.IndexOptions{
				ExpireAfterSeconds: &ttl,
			},
		}, opts)
	if err != nil {
		log.Fatalln(err)
	}

	sessionIndex, err := db.Collection("sessions").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{"lastActive": 1},
			Options: &options.IndexOptions{
				ExpireAfterSeconds: &ttl,
			},
		}, opts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Indexes:", emailIndex, tokenIndex, sessionIndex)

	fmt.Println("Connected to DATABASE:", dbName)

	return db
}
