package db

import (
	"context"
	"github.com/mustthink/go-test-api/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ConnClient(url string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}

func CheckDB(client *mongo.Client) (bool, error) {
	collection := client.Database("test").Collection("transactions")

	var result types.Transaction
	filter := bson.D{}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != mongo.ErrNoDocuments {
		return true, err
	} else if err == nil {
		return true, nil
	}

	return false, nil
}
