package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DatabaseName             = "inventory"
	IngredientCollectionName = "ingredient"
)

// getting database collections
func GetIngredientCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(DatabaseName).Collection(IngredientCollectionName)
}
