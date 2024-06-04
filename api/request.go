package api

type IngredientRequest struct {
	ID       string  `json:"id" validate:"required"`
	Name     string  `json:"name" validate:"omitempty"`
	Quantity float64 `json:"quantity" validate:"required,min=0.1"`
	Units    string  `json:"units" validate:"oneof=i is cs tbsp tsp g kg"`
}
