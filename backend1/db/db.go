package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Db() *mongo.Client {
	// clientOptions := options.Client().ApplyURI("mongodb+srv://ryanmolley05:i71h3F1Bxt1UkLVf@primaryexam.yapdqkj.mongodb.net/?retryWrites=true&w=majority")
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}
