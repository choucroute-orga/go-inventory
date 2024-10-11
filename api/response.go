package api

import (
	"inventory/db"
	"time"
)

const (
	LiveStatus     = "OK"
	ReadyStatus    = "READY"
	NotReadyStatus = "NOT READY"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func NewHealthResponse(status string) *HealthResponse {
	return &HealthResponse{
		Status: status,
	}
}

type IngredientResponse struct {
	ID        string      `json:"id"` // ingredientID
	Name      string      `json:"name"`
	Amount    float64     `json:"amount"`
	Unit      db.UnitType `json:"unit"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

func NewIngredientResponse(ingredient *db.UserInventory) *IngredientResponse {
	return &IngredientResponse{
		ID:        ingredient.IngredientID,
		Name:      ingredient.Name,
		Amount:    ingredient.Quantity,
		Unit:      ingredient.Unit,
		CreatedAt: ingredient.CreatedAt,
		UpdatedAt: ingredient.UpdatedAt,
	}
}

func NewAllIngredientsResponse(ingredients *[]db.UserInventory) *[]IngredientResponse {
	ingredientsResponse := []IngredientResponse{}

	for _, ingredient := range *ingredients {
		ingredientsResponse = append(ingredientsResponse, IngredientResponse{
			ID:        ingredient.IngredientID,
			Name:      ingredient.Name,
			Amount:    ingredient.Quantity,
			Unit:      ingredient.Unit,
			CreatedAt: ingredient.CreatedAt,
			UpdatedAt: ingredient.UpdatedAt,
		})
	}
	return &ingredientsResponse
}
