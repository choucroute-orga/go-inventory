package api

import (
	"encoding/json"
	"inventory/db"
	"inventory/messages"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (api *ApiHandler) createShoppingList(l *logrus.Entry, recipe messages.AddRecipe) *[]messages.AddIngredient {

	ingredientsShoppingList := []messages.AddIngredient{}

	for _, ingredient := range recipe.Ingredients {
		// Check if the ingredient is in the inventory
		// If not, send the ingredient to the shopping list
		id, err := primitive.ObjectIDFromHex(ingredient.ID)
		if err != nil {
			logger.WithError(err).Info("Failed to convert the ID")
			continue
		}
		ingredientInInventory, err := db.FindById(l, api.mongo, id)
		l := l.WithFields(logrus.Fields{
			"id":         ingredient.ID,
			"ingredient": ingredientInInventory.Name,
			"amount":     ingredient.Amount,
			"unit":       ingredientInInventory.Units})

		if err != nil {
			ingredientsShoppingList = append(ingredientsShoppingList, ingredient)
		} else if ingredientInInventory.Units != ingredient.Unit {
			ingredientsShoppingList = append(ingredientsShoppingList, ingredient)
		} else if ingredientInInventory.Quantity < ingredient.Amount {
			// Calculate the amount to buy
			ingredient.Amount = ingredient.Amount - ingredientInInventory.Quantity
			l = l.WithField("amount", ingredient.Amount)
			ingredientsShoppingList = append(ingredientsShoppingList, ingredient)
		} else {
			l = l.WithField("amount", 0)
		}
		l.Debug("Ingredient treated")
	}

	return &ingredientsShoppingList
}

func (api *ApiHandler) ConsumesMessages() {
	for {
		api.consumesMessages()
		time.Sleep(time.Microsecond)
	}
}

// This function is used to consume messages from the inventory queue
// It retrieves the request of adding a recipe into the shopping list
// It check if ingredients are in the inventory and send the ingredient to be added to the shopping list
func (api *ApiHandler) consumesMessages() {

	l := logger.WithField("context", "consumesMessages")

	q := messages.GetInventoryQueue(api.amqp)
	if q == nil {
		return
	}

	ch, err := api.amqp.Channel()
	if err != nil {
		logger.WithError(err).Error("Failed to open a channel")
	}

	msgs, err := ch.Consume(
		q.Name,      // queue
		"inventory", // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	defer ch.Close()

	if err != nil {
		logger.WithError(err).Error("Failed to register a consumer")
	}

	for d := range msgs {

		logger.WithField("message", string(d.Body)).Info("Received a message")
		var recipe messages.AddRecipe
		err := json.Unmarshal(d.Body, &recipe)
		if err != nil {
			logger.WithError(err).Error("Failed to unmarshal the message")
		}
		err = api.validation.Validate.Struct(recipe)
		if err != nil {
			logger.WithError(err).Error("Failed to validate the message")
			break
		}
		// Check all ingredients and the amount
		ingredientsShoppingList := api.createShoppingList(l, recipe)

		recipe.Ingredients = *ingredientsShoppingList
		// Publish the ingredients to the shopping list

		// Convert the error into a JSON
		jsonMessage, err := json.Marshal(recipe)
		if err != nil {
			logger.WithError(err).Error("Failed to marshal the recipe")
		}

		qSL := messages.GetShoppingListQueue(api.amqp)

		// Publish in the inventory-error queue the error
		err = ch.Publish(
			"",       // exchange
			qSL.Name, // routing key
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				ContentType: "application/json", // JSON
				Body:        jsonMessage,        // The error
			})

		if err != nil {
			logger.WithError(err).Error("Failed to publish the message")
		}
		logger.WithField("message", string(jsonMessage)).Info("Published the shopping list message")

	}

}
