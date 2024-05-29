package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("inventory").Collection(collectionName)
}
