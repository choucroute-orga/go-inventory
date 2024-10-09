package api

import "inventory/db"

const (
	LiveStatus     = "OK"
	ReadyStatus    = "READY"
	NotReadyStatus = "NOT READY"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type Quantity struct {
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"` // i is cs tbsp tsp g kg
}

type IngredientResponse struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Quantities []Quantity `json:"quantities"`
}

func NewHealthResponse(status string) *HealthResponse {
	return &HealthResponse{
		Status: status,
	}
}

func NewIngredientResponse(ingredients *[]db.Ingredient) *IngredientResponse {
	id := (*ingredients)[0].ID.Hex()
	name := (*ingredients)[0].Name
	quantites := make([]Quantity, len(*ingredients))

	for i, ingredient := range *ingredients {
		quantites[i] = Quantity{
			Amount: ingredient.Quantity,
			Unit:   ingredient.Units,
		}
	}

	ingredient := IngredientResponse{
		ID:         id,
		Name:       name,
		Quantities: quantites,
	}
	return &ingredient
}
