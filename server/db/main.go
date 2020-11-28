package db

import (
	"clipboard/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Db() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Connect to //MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

func InitDb() {
	Db()
	CreateIndex(models.CollectionUser, "username", true)
	CreateIndex(models.CollectionUser, "created_at", false)

}
