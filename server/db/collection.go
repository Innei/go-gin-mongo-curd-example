package db

import (
	"clipboard/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Database = Db().Database("clipboard")
var ClipCollection = Database.Collection(models.CollectionClip)
var UserCollection = Database.Collection(models.CollectionUser)

// CreateIndex - creates an index for a specific field in a collection
func CreateIndex(collectionName string, field string, unique bool) bool {

	// 1. Lets define the keys for the index we want to create
	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique),
	}

	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Connect to the database and access the collection
	collection := Database.Collection(collectionName)

	// 4. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		return false
	}

	// 6. All went well, we return true
	return true
}
