package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ingredient struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	IngredientID string             `bson:"ingredient_id" json:"ingredient_id"`
	Name         string             `bson:"name" json:"name"`
	Quantity     float64            `bson:"quantity" json:"amount"`
	Units        string             `bson:"units" json:"unit"`
}
