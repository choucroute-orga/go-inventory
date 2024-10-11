package api

import (
	"encoding/json"
	"inventory/db"
	"inventory/messages"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func (api *ApiHandler) createShoppingList(l *logrus.Entry, recipe messages.AddRecipe) *[]messages.NeededIngredientShoppingList {
	l = logger.WithField("context", "createShoppingList")
	var shoppingList []messages.NeededIngredientShoppingList

	// Get user's current inventory
	userInventory, err := db.GetAll(l, api.mongo, recipe.UserID)
	if err != nil {
		l.WithError(err).Error("Failed to fetch user inventory")
		return &shoppingList
	}

	// Create a map for quick inventory lookup
	inventoryMap := make(map[string]db.UserInventory)
	for _, item := range userInventory {
		inventoryMap[item.IngredientID] = item
	}

	// Process each ingredient in the recipe
	for _, recipeIngredient := range recipe.Ingredients {
		// Convert recipe ingredient quantity to base unit
		res, err := ConvertToBaseUnitFromRequest(recipeIngredient.Amount, recipeIngredient.Unit)
		recipeQtyBase, baseUnit := res.Quantity, res.Unit
		if err != nil {
			l.WithError(err).WithFields(logrus.Fields{
				"ingredientId": recipeIngredient.ID,
				"unit":         recipeIngredient.Unit,
			}).Error("Failed to convert recipe ingredient to base unit")
			continue
		}

		// Check if user has this ingredient
		userItem, exists := inventoryMap[recipeIngredient.ID]
		if !exists {
			// User doesn't have this ingredient at all - add full amount converted to shopping list
			messageConv := roundBaseUnit(res)
			finalQty, finalUnit := messageConv.Quantity, messageConv.Unit
			shoppingList = append(shoppingList, messages.NeededIngredientShoppingList{
				ID:     recipeIngredient.ID,
				Amount: finalQty,
				Unit:   finalUnit,
			})
			continue
		}

		// Convert user's inventory quantity to the same base unit
		res, err = ConvertToBaseUnit(userItem.Quantity, userItem.Unit)
		userQtyBase, userBaseUnit := res.Quantity, res.Unit
		if err != nil {
			l.WithError(err).WithFields(logrus.Fields{
				"ingredientId": recipeIngredient.ID,
				"unit":         userItem.Unit,
			}).Error("Failed to convert user inventory to base unit")
			continue
		}

		// Check if units are compatible
		if (baseUnit == db.UnitMl && userBaseUnit != db.UnitMl) ||
			(baseUnit == db.UnitG && userBaseUnit != db.UnitG) ||
			(baseUnit == db.UnitItem && userBaseUnit != db.UnitItem) {
			l.WithFields(logrus.Fields{
				"recipeUnit": baseUnit,
				"userUnit":   userBaseUnit,
			}).Error("Incompatible units")
			continue
		}

		// Calculate if more is needed
		if userQtyBase < recipeQtyBase {
			neededQty := recipeQtyBase - userQtyBase

			conversionMessage := roundBaseUnit(ConversionResult{neededQty, baseUnit})
			finalQty, finalUnit := conversionMessage.Quantity, conversionMessage.Unit

			shoppingList = append(shoppingList, messages.NeededIngredientShoppingList{
				ID:     recipeIngredient.ID,
				Amount: finalQty,
				Unit:   finalUnit,
			})
		}
	}

	return &shoppingList
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

		sendRecipe := messages.SendRecipe{
			ID:          recipe.ID,
			UserID:      recipe.UserID,
			Ingredients: *ingredientsShoppingList,
		}
		// Publish the ingredients to the shopping list

		// Convert the error into a JSON
		jsonMessage, err := json.Marshal(sendRecipe)
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
