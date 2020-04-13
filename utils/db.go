package utils

import (
	"context"
	"fmt"
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
		log.Fatal("Missing dev db connection string")
	}

	dbName, nameExists := os.LookupEnv("DB_NAME")
	if !nameExists {
		log.Fatal("Missing dev db name")
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

	fmt.Println("Connected to DATABASE:", dbName)

	return db
}
