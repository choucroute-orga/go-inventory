package api

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory/db"
	"inventory/messages"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func (api *ApiHandler) ConsumeAddRecipeMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping ShoppingListAddRecipe message consumption")
			return
		default:
			api.consumeAddRecipeMessages(ctx)
			time.Sleep(time.Second)
		}
	}
}

func (api *ApiHandler) ConsumeAddIngredientMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping ShoppingListAddIngredient message consumption")
			return
		default:
			api.consumeAddIngredientMessage(ctx)
			time.Sleep(time.Second)
		}
	}
}

func (api *ApiHandler) consumeAddIngredientMessage(ctx context.Context) {
	l := logger.WithField("context", "consumeAddIngredientMessage")

	q := messages.GetInventoryIngredientQueue(api.amqp)
	if q == nil {
		return
	}
	ch, _ := messages.OpenChannel(api.amqp)
	if ch == nil {
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name,      // queue
		"inventory", // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	l = l.WithFields(logrus.Fields{
		"queue": q.Name})

	if err != nil {
		logger.WithError(err).Error("Failed to register a consumer")
		return
	}

	l.Info("Started consuming messages")

	for {
		select {
		case <-ctx.Done():
			l.Info("Shutting down consumer")
			return
		case msg, ok := <-msgs:
			if !ok {
				l.Warn("Message channel closed")
				return
			}
			if err := api.processMessageAddIngredient(l, msg); err != nil {
				l.WithError(err).Error("Failed to process message")
				// Implement retry logic or move to dead-letter queue
			} else {
				msg.Ack(false)
			}
		}
	}
}

func (api *ApiHandler) processMessageAddIngredient(l *logrus.Entry, msg amqp.Delivery) error {

	l = l.WithField("function", "processMessageAddIngredient")
	var ingredient messages.IngredientInventory
	if err := json.Unmarshal(msg.Body, &ingredient); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	if err := api.validation.Validate.Struct(ingredient); err != nil {
		return fmt.Errorf("failed to validate message: %w", err)
	}
	ingredientShoppingList, err := api.createIngredientShoppingList(l, ingredient)
	if err != nil {
		return fmt.Errorf("failed to create ingredient shopping list: %w", err)
	}
	return api.publishIngredientShoppingList(l, *ingredientShoppingList)
}

func (api *ApiHandler) createIngredientShoppingList(l *logrus.Entry, ingredient messages.IngredientInventory) (*messages.IngredientShoppingList, error) {
	l = l.WithField("function", "createIngredientShoppingList")
	neededIngredient, err := api.processIngredient(l, ingredient.UserID, ingredient.Ingredient)
	if err != nil {
		l.WithError(err).WithField("ingredientId", ingredient.Ingredient.ID).Error("Failed to process ingredient")
		return nil, err
	}
	ingredientShoppingList := messages.IngredientShoppingList{
		NeededIngredientShoppingList: *neededIngredient,
		UserID:                       ingredient.UserID,
	}
	return &ingredientShoppingList, nil
}

func (api *ApiHandler) publishIngredientShoppingList(l *logrus.Entry, ingredient messages.IngredientShoppingList) error {
	l = l.WithField("function", "publishIngredientShoppingList")

	ch, err := messages.OpenChannel(api.amqp)
	if ch == nil {
		return err
	}

	jsonMessage, err := json.Marshal(ingredient)
	if err != nil {
		l.WithError(err).Error("Failed to marshal the ingredient")
	}

	qSL := messages.GetShoppingListIngredientQueue(api.amqp)

	err = ch.Publish(
		"",       // exchange
		qSL.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonMessage,
		})

	if err != nil {
		logger.WithError(err).Error("Failed to publish the message")
	}
	logger.WithField("message", string(jsonMessage)).Info("Published the Ingredient shopping list message")
	return nil
}

// This function is used to consume messages from the inventory queue
// It retrieves the request of adding a recipe into the shopping list
// It check if ingredients are in the inventory and send the ingredient to be added to the shopping list
func (api *ApiHandler) consumeAddRecipeMessages(ctx context.Context) {

	l := logger.WithField("context", "consumeAddRecipeMessages")

	q := messages.GetInventoryRecipeQueue(api.amqp)
	if q == nil {
		return
	}

	ch, _ := messages.OpenChannel(api.amqp)
	if ch == nil {
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name,      // queue
		"inventory", // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	l = l.WithFields(logrus.Fields{
		"queue": q.Name})

	if err != nil {
		logger.WithError(err).Error("Failed to register a consumer")
		return
	}

	l.Info("Started consuming messages")

	for {
		select {
		case <-ctx.Done():
			l.Info("Shutting down consumer")
			return
		case msg, ok := <-msgs:
			if !ok {
				l.Warn("Message channel closed")
				return
			}
			if err := api.processMessage(l, msg); err != nil {
				l.WithError(err).Error("Failed to process message")
				// Implement retry logic or move to dead-letter queue
			} else {
				msg.Ack(false)
			}
		}
	}

}

func (api *ApiHandler) processMessage(l *logrus.Entry, msg amqp.Delivery) error {
	l = l.WithField("function", "processMessage")
	var recipe messages.AddRecipe
	if err := json.Unmarshal(msg.Body, &recipe); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if err := api.validation.Validate.Struct(recipe); err != nil {
		return fmt.Errorf("failed to validate message: %w", err)
	}

	ingredientsShoppingList := api.createShoppingList(l, recipe)

	sendRecipe := messages.SendRecipe{
		ID:          recipe.ID,
		UserID:      recipe.UserID,
		Ingredients: *ingredientsShoppingList,
	}

	return api.publishShoppingList(l, sendRecipe)
}

func (api *ApiHandler) createShoppingList(l *logrus.Entry, recipe messages.AddRecipe) *[]messages.NeededIngredientShoppingList {
	l = logger.WithField("context", "createShoppingList")
	var shoppingList []messages.NeededIngredientShoppingList

	// Process each ingredient in the recipe
	for _, recipeIngredient := range recipe.Ingredients {
		neededIngredient, err := api.processIngredient(l, recipe.UserID, recipeIngredient)
		if err != nil {
			l.WithError(err).WithField("ingredientId", recipeIngredient.ID).Error("Failed to process ingredient")
			continue
		}
		if neededIngredient != nil {
			shoppingList = append(shoppingList, *neededIngredient)
		}

	}

	return &shoppingList
}

func (api *ApiHandler) processIngredient(l *logrus.Entry, userID string, ingredient messages.Ingredient) (*messages.NeededIngredientShoppingList, error) {
	// Convert ingredient quantity to base unit
	res, err := ConvertToBaseUnitFromRequest(ingredient.Amount, ingredient.Unit)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ingredient to base unit: %w", err)
	}
	ingredientQtyBase, baseUnit := res.Quantity, res.Unit

	// Get user's current inventory for this ingredient
	userItem, err := db.GetOne(l, api.mongo, userID, ingredient.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User doesn't have this ingredient at all - return full amount converted
			messageConv := roundBaseUnit(res)
			return &messages.NeededIngredientShoppingList{
				ID:     ingredient.ID,
				Amount: messageConv.Quantity,
				Unit:   messageConv.Unit,
			}, nil
		}
		return nil, fmt.Errorf("failed to fetch user inventory: %w", err)
	}

	// Convert user's inventory quantity to the same base unit
	userRes, err := ConvertToBaseUnit(userItem.Quantity, userItem.Unit)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user inventory to base unit: %w", err)
	}
	userQtyBase, userBaseUnit := userRes.Quantity, userRes.Unit

	// Check if units are compatible
	if (baseUnit == db.UnitMl && userBaseUnit != db.UnitMl) ||
		(baseUnit == db.UnitG && userBaseUnit != db.UnitG) ||
		(baseUnit == db.UnitItem && userBaseUnit != db.UnitItem) {
		return nil, fmt.Errorf("incompatible units: recipe %s, user %s", baseUnit, userBaseUnit)
	}

	// Calculate if more is needed
	if userQtyBase < ingredientQtyBase {
		neededQty := ingredientQtyBase - userQtyBase
		conversionMessage := roundBaseUnit(ConversionResult{neededQty, baseUnit})
		return &messages.NeededIngredientShoppingList{
			ID:     ingredient.ID,
			Amount: conversionMessage.Quantity,
			Unit:   conversionMessage.Unit,
		}, nil
	}

	return nil, nil // No additional quantity needed
}

func (api *ApiHandler) publishShoppingList(l *logrus.Entry, sendRecipe messages.SendRecipe) error {
	l = l.WithField("function", "publishShoppingList")

	ch, err := messages.OpenChannel(api.amqp)
	if ch == nil {
		return err
	}

	// TODO Convert the error into a JSON message and send to a queue where we can read it
	jsonMessage, err := json.Marshal(sendRecipe)
	if err != nil {
		l.WithError(err).Error("Failed to marshal the recipe")
	}

	qSL := messages.GetShoppingListRecipeQueue(api.amqp)

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
	return nil
}
