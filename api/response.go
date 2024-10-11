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

// func NewIngredientResponse(ingredients *[]db.Ingredient) *IngredientResponse {

// 	var ingredient IngredientResponse

// 	if len(*ingredients) == 0 {
// 		return nil
// 	}

// 	id := (*ingredients)[0].ID.Hex()
// 	name := (*ingredients)[0].Name
// 	quantites := make([]Quantity, len(*ingredients))

// 	for i, ingredient := range *ingredients {
// 		quantites[i] = Quantity{
// 			Amount: ingredient.Quantity,
// 			Unit:   ingredient.Units,
// 		}
// 	}

// 	ingredient = IngredientResponse{
// 		ID:         id,
// 		Name:       name,
// 		Quantities: quantites,
// 	}
// 	return &ingredient
// }

// func NewAllIngredientsResponse(ingredients *[]db.Ingredient) *[]IngredientResponse {
// 	// We append ingredient by ingredientID together
// 	ingredientsResponse := []IngredientResponse{}

// 	for _, ingredient := range *ingredients {
// 		ingredientFound := false
// 		for i, ingredientResponse := range ingredientsResponse {
// 			if ingredientResponse.ID == ingredient.IngredientID {
// 				ingredientsResponse[i].Quantities = append(ingredientsResponse[i].Quantities, Quantity{
// 					Amount: ingredient.Quantity,
// 					Unit:   ingredient.Units,
// 				})
// 				ingredientFound = true
// 			}
// 		}
// 		if !ingredientFound {
// 			ingredientsResponse = append(ingredientsResponse, IngredientResponse{
// 				ID:   ingredient.IngredientID,
// 				Name: ingredient.Name,
// 				Quantities: []Quantity{
// 					{
// 						Amount: ingredient.Quantity,
// 						Unit:   ingredient.Units,
// 					},
// 				},
// 			})
// 		}
// 	}
// 	return &ingredientsResponse
// }
