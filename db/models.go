package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ingredient struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
	Unit string             `bson:"unit" json:"unit"`
}
