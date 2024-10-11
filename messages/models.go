package messages

type UnitRequest string

const (
	UnitItem  UnitRequest = "i"
	UnitItems UnitRequest = "is"
	UnitG     UnitRequest = "g"
	UnitKg    UnitRequest = "kg"
	UnitMl    UnitRequest = "ml"
	UnitL     UnitRequest = "l"
	UnitTsp   UnitRequest = "tsp"
	UnitTbsp  UnitRequest = "tbsp"
	UnitCs    UnitRequest = "cs"
)

type Ingredient struct {
	ID     string      `param:"id" json:"id" validate:"required"`
	Amount float64     `json:"amount" validate:"required,min=0"`
	Unit   UnitRequest `json:"unit" validate:"oneof=i is cs tbsp tsp g kg ml l"`
}

// This is the message that is received from the Gateway
type AddRecipe struct {
	ID string `json:"id" validate:"required"`
	// TODO Add userID in the other MS
	UserID      string       `json:"userId" validate:"required"`
	Ingredients []Ingredient `json:"ingredients" validate:"required,dive,required"`
}

type SendRecipe struct {
	ID string `json:"id" validate:"required"`
	// TODO Add userID in the other MS
	UserID      string                         `json:"userId" validate:"required"`
	Ingredients []NeededIngredientShoppingList `json:"ingredients" validate:"required,dive,required"`
}

type UnitMessage string

const (
	UnitMessageItem UnitMessage = "i"
	UnitMessageG    UnitMessage = "g"
	UnitMessageKg   UnitMessage = "kg"
	UnitMessageMl   UnitMessage = "ml"
	UnitMessageL    UnitMessage = "l"
)

type ConversionMessageResult struct {
	Quantity float64     `json:"quantity" validate:"required,min=0"`
	Unit     UnitMessage `json:"unit" validate:"required,oneof=i g ml l"`
}

// This is the message that is sent to the Shopping List MS for the needed ingredients
type NeededIngredientShoppingList struct {
	ID     string      `json:"id" validate:"required"`
	Amount float64     `json:"amount" validate:"required,min=0"`
	Unit   UnitMessage `json:"unit" validate:"required,oneof=i g ml"`
}
