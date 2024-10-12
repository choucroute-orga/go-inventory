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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	ctx, span := api.tracer.Start(ctx, "consumeAddIngredientMessage")
	defer span.End()
	l := logger.WithField("context", "consumeAddIngredientMessage")
	ch, err := messages.OpenChannel(api.amqp)
	if err != nil {
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		messages.InventoryAddIngredientShoppingList, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		logger.WithError(err).Errorf("Failed to declare queue %s", q.Name)
		return
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
			messageCtx, messageSpan := api.tracer.Start(ctx, "processShoppingListAddRecipeMessage")
			startTime := time.Now()

			retryCount := 0
			maxRetries := 3
			var processErr error
			for retryCount < maxRetries {
				processErr = api.processAddIngredientMessage(messageCtx, l, msg)
				if processErr == nil {
					break
				}
				retryCount++
				l.WithError(processErr).WithField("retry", retryCount).Warn("Retrying message processing")
				time.Sleep(time.Second * time.Duration(retryCount)) // Exponential backoff
			}

			duration := time.Since(startTime)
			processStatus := "success"
			if processErr != nil {
				processStatus = "failure"
			}
			messageSpan.SetAttributes(
				attribute.Int("retries", retryCount),
				attribute.String("status", processStatus),
				attribute.Int64("duration_ms", duration.Milliseconds()),
			)
			messageSpan.End()

			if processErr != nil {
				l.WithError(processErr).Error("Failed to process message after max retries")
				// Send to dead-letter queue
				err := ch.Publish(
					"",                           // exchange
					messages.DeadLetterQueueName, // routing key
					false,                        // mandatory
					false,                        // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        msg.Body,
						Headers: amqp.Table{
							"x-original-queue": messages.AddIngredientShoppingList,
							"x-error":          processErr.Error(),
						},
					},
				)
				if err != nil {
					l.WithError(err).Error("Failed to send message to dead-letter queue")
				}
			}

			msg.Ack(false)
		}
	}
}

func (api *ApiHandler) processAddIngredientMessage(ctx context.Context, l *logrus.Entry, msg amqp.Delivery) error {
	ctx, span := api.tracer.Start(ctx, "processAddIngredientMessage")
	defer span.End()
	l = l.WithField("function", "processMessageAddIngredient")
	var ingredient messages.IngredientInventory
	if err := json.Unmarshal(msg.Body, &ingredient); err != nil {
		span.SetStatus(codes.Error, "Failed to unmarshal message")
		span.RecordError(err)
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	span.SetAttributes(
		attribute.String("ingredient.id", ingredient.Ingredient.ID),
		attribute.String("ingredient.user_id", ingredient.UserID),
		attribute.Float64("ingredient.amount", ingredient.Ingredient.Amount),
		attribute.String("ingredient.unit", string(ingredient.Ingredient.Unit)),
	)

	if err := api.validation.Validate.Struct(ingredient); err != nil {
		span.SetStatus(codes.Error, "Failed to validate message")
		span.RecordError(err)
		return fmt.Errorf("failed to validate message: %w", err)
	}
	ingredientShoppingList, err := api.createIngredientShoppingList(ctx, l, ingredient)
	if err != nil {
		return fmt.Errorf("failed to create ingredient shopping list: %w", err)
	}
	return api.publishIngredientShoppingList(l, *ingredientShoppingList)
}

func (api *ApiHandler) createIngredientShoppingList(ctx context.Context, l *logrus.Entry, ingredient messages.IngredientInventory) (*messages.IngredientShoppingList, error) {
	ctx, span := api.tracer.Start(ctx, "createIngredientShoppingList")
	defer span.End()
	l = l.WithField("function", "createIngredientShoppingList")
	neededIngredient, err := api.processIngredient(ctx, l, ingredient.UserID, ingredient.Ingredient)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to process ingredient")
		span.RecordError(err)
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
	if err != nil {
		return err
	}

	jsonMessage, err := json.Marshal(ingredient)
	if err != nil {
		l.WithError(err).Error("Failed to marshal the ingredient")
	}

	qSL, err := ch.QueueDeclare(
		messages.AddIngredientShoppingList, // name
		true,                               // durable
		false,                              // delete when unused
		false,                              // exclusive
		false,                              // no-wait
		nil,                                // arguments
	)

	if err != nil {
		logger.WithError(err).Error("Failed to declare a queue")
	}

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
	
	l := logger.WithContext(ctx).WithField("function", "consumeAddRecipeMessages")

	ch, err := messages.OpenChannel(api.amqp)
	if err != nil {
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		messages.InventoryAddRecipesShoppingList, // name
		true,                                     // durable
		false,                                    // delete when unused
		false,                                    // exclusive
		false,                                    // no-wait
		nil,                                      // arguments
	)
	if err != nil {
		logger.WithError(err).Errorf("Failed to declare queue %v", q.Name)
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
			messageCtx, messageSpan := api.tracer.Start(ctx, "processShoppingListAddRecipeMessage")
			startTime := time.Now()

			retryCount := 0
			maxRetries := 3
			var processErr error
			for retryCount < maxRetries {
				processErr = api.processShoppingListAddRecipeMessage(messageCtx, l, msg)
				if processErr == nil {
					break
				}
				retryCount++
				l.WithError(processErr).WithField("retry", retryCount).Warn("Retrying message processing")
				time.Sleep(time.Second * time.Duration(retryCount)) // Exponential backoff
			}

			duration := time.Since(startTime)
			processStatus := "success"
			if processErr != nil {
				processStatus = "failure"
			}
			messageSpan.SetAttributes(
				attribute.Int("retries", retryCount),
				attribute.String("status", processStatus),
				attribute.Int64("duration_ms", duration.Milliseconds()),
			)
			messageSpan.End()

			if processErr != nil {
				l.WithError(processErr).Error("Failed to process message after max retries")
				// Send to dead-letter queue
				err := ch.Publish(
					"",                           // exchange
					messages.DeadLetterQueueName, // routing key
					false,                        // mandatory
					false,                        // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        msg.Body,
						Headers: amqp.Table{
							"x-original-queue": messages.AddRecipesShoppingList,
							"x-error":          processErr.Error(),
						},
					},
				)
				if err != nil {
					l.WithError(err).Error("Failed to send message to dead-letter queue")
				}
			}

			msg.Ack(false)
		}
	}

}

func (api *ApiHandler) processShoppingListAddRecipeMessage(ctx context.Context, l *logrus.Entry, msg amqp.Delivery) error {
	ctx, span := api.tracer.Start(ctx, "processShoppingListAddRecipeMessage")
	defer span.End()
	l = l.WithContext(ctx).WithField("function", "processShoppingListAddRecipeMessage")
	var recipe messages.AddRecipe
	if err := json.Unmarshal(msg.Body, &recipe); err != nil {
		span.SetStatus(codes.Error, "Failed to unmarshal message")
		span.RecordError(err)
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	span.SetAttributes(
		attribute.String("recipe.id", recipe.ID),
		attribute.String("recipe.user_id", recipe.UserID),
		attribute.Int("recipe.ingredients", len(recipe.Ingredients)),
	)

	if err := api.validation.Validate.Struct(recipe); err != nil {
		span.SetStatus(codes.Error, "Failed to validate message")
		span.RecordError(err)
		return fmt.Errorf("failed to validate message: %w", err)
	}
	ingredientsCtx, ingredientsSpan := api.tracer.Start(ctx, "createShoppingList")
	ingredientsShoppingList := api.createShoppingList(ingredientsCtx, l, recipe)
	ingredientsSpan.SetAttributes(attribute.Int("ingredients_count", len(*ingredientsShoppingList)))
	ingredientsSpan.End()

	sendRecipe := messages.SendRecipe{
		ID:          recipe.ID,
		UserID:      recipe.UserID,
		Ingredients: *ingredientsShoppingList,
	}

	publishCtx, publishSpan := api.tracer.Start(ctx, "publishShoppingList")
	l = l.WithContext(publishCtx)
	err := api.publishShoppingList(l, sendRecipe)
	if err != nil {
		publishSpan.SetStatus(codes.Error, "Failed to publish shopping list")
		publishSpan.RecordError(err)
	}
	publishSpan.End()
	return err
}

func (api *ApiHandler) createShoppingList(ctx context.Context, l *logrus.Entry, recipe messages.AddRecipe) *[]messages.NeededIngredientShoppingList {
	l = logger.WithField("context", "createShoppingList")
	var shoppingList []messages.NeededIngredientShoppingList

	// Process each ingredient in the recipe
	for _, recipeIngredient := range recipe.Ingredients {
		neededIngredient, err := api.processIngredient(ctx, l, recipe.UserID, recipeIngredient)
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

func (api *ApiHandler) processIngredient(ctx context.Context, l *logrus.Entry, userID string, ingredient messages.Ingredient) (*messages.NeededIngredientShoppingList, error) {
	// Convert ingredient quantity to base unit
	msgCtx, span := api.tracer.Start(ctx, "processIngredient")
	l = l.WithContext(msgCtx).WithField("function", "processIngredient")
	defer span.End()
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
			span.SetAttributes(
				attribute.String("ingredient.id", ingredient.ID),
				attribute.Float64("ingredient.amount", messageConv.Quantity),
				attribute.String("ingredient.unit", string(messageConv.Unit)),
			)
			return &messages.NeededIngredientShoppingList{
				ID:     ingredient.ID,
				Amount: messageConv.Quantity,
				Unit:   messageConv.Unit,
			}, nil
		}
		span.SetStatus(codes.Error, "Failed to fetch user inventory")
		span.RecordError(err)
		return nil, fmt.Errorf("failed to fetch user inventory: %w", err)
	}

	// Convert user's inventory quantity to the same base unit
	userRes, err := ConvertToBaseUnit(userItem.Quantity, userItem.Unit)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to convert user inventory to base unit")
		span.RecordError(err)
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
		span.SetAttributes(
			attribute.String("ingredient.id", ingredient.ID),
			attribute.Float64("ingredient.amount", conversionMessage.Quantity),
			attribute.String("ingredient.unit", string(conversionMessage.Unit)),
		)
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
	if err != nil {
		return err
	}

	// TODO Convert the error into a JSON message and send to a queue where we can read it
	jsonMessage, err := json.Marshal(sendRecipe)
	if err != nil {
		l.WithError(err).Error("Failed to marshal the recipe")
		return err
	}

	qSL, err := ch.QueueDeclare(
		messages.AddRecipesShoppingList, // name
		true,                            // durable
		false,                           // delete when unused
		false,                           // exclusive
		false,                           // no-wait
		nil,                             // arguments
	)

	if err != nil {
		logger.WithError(err).Error("Failed to declare a queue")
		return err
	}

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
