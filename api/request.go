package api

import (
	"inventory/db"
	"inventory/messages"
)

type PostIngredientRequest struct {
	ID     string               `json:"id" validate:"required"`
	UserID string               `json:"userId" validate:"required"`
	Name   string               `json:"name" validate:"omitempty"`
	Amount float64              `json:"amount" validate:"required,min=0.1"`
	Unit   messages.UnitRequest `json:"unit" validate:"oneof=i is cs tbsp tsp g kg ml l"`
}

// It's the same than the post, just that the ID is in the URL
type PutIngredientRequest struct {
	ID     string               `param:"id" validate:"required"`
	UserID string               `json:"userId" validate:"required"`
	Name   string               `json:"name" validate:"omitempty"`
	Amount float64              `json:"amount" validate:"required,min=0.1"`
	Unit   messages.UnitRequest `json:"unit" validate:"oneof=i is cs tbsp tsp g kg ml l"`
}

type DeleteIngredientRequest struct {
	ID     string `param:"id" validate:"required"`
	UserID string `param:"userId" validate:"required"`
}

func NewIngredientInventoryFromPut(ingredient *PutIngredientRequest) *db.UserInventory {

	// Convert the ingredient to the base unit
	res, _ := ConvertToBaseUnitFromRequest(ingredient.Amount, ingredient.Unit)

	return &db.UserInventory{
		IngredientID: ingredient.ID,
		UserID:       ingredient.UserID,
		Name:         ingredient.Name,
		Quantity:     res.Quantity,
		Unit:         res.Unit,
	}
}

func NewIngredientInventory(ingredient *PostIngredientRequest) *db.UserInventory {

	// Convert the ingredient to the base unit
	res, _ := ConvertToBaseUnitFromRequest(ingredient.Amount, ingredient.Unit)

	return &db.UserInventory{
		IngredientID: ingredient.ID,
		UserID:       ingredient.UserID,
		Name:         ingredient.Name,
		Quantity:     res.Quantity,
		Unit:         res.Unit,
	}
}
