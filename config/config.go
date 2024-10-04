package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client                *mongo.Client
	TransactionCollection *mongo.Collection
)

func InitMongo() {
	password := os.Getenv("MONGODB_PASSWORD")
	dbName := os.Getenv("MONGODB_DBNAME")
    
	if password == "" || dbName == "" {
		log.Fatal("MONGODB_PASSWORD or MONGODB_DBNAME environment variable not set")
	}

	uri := fmt.Sprintf("mongodb+srv://trocup:%s@pli-etna.tpeqiyq.mongodb.net/?retryWrites=true&w=majority&appName=PLI-ETNA", password)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err := Client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	InitTransactionCollection(dbName)
}

func InitTransactionCollection(dbName string) {
	TransactionCollection = Client.Database(dbName).Collection("transaction")
}

func CleanUpTestDatabase(dbName string) {
	if err := Client.Database(dbName).Drop(context.TODO()); err != nil {
		log.Fatalf("Failed to clean up test database: %v", err)
	}
	fmt.Println("Test database cleaned up successfully.")
}
