package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UnitType represents the measurement unit for ingredients
type UnitType string

const (
	UnitItem UnitType = "i"
	UnitG    UnitType = "g"
	UnitMl   UnitType = "ml"
)

// Ingredient represents an ingredient in a user's inventory
type UserInventory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	UserID       string             `bson:"userId" json:"userId"`
	IngredientID string             `bson:"ingredientId" json:"ingredientId"` // Reference to the ingredient in the ingredient service
	Quantity     float64            `bson:"quantity" json:"amount"`
	Unit         UnitType           `bson:"unit" json:"unit"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Create a String method for the UnitType
func (u UnitType) String() string {
	return string(u)
}
