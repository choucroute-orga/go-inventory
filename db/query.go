package db

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var loger = logrus.WithFields(logrus.Fields{
	"context": "db/query",
})

func NewID() primitive.ObjectID {
	return primitive.NewObjectIDFromTimestamp(time.Now())
}

// GetAll retrieves all ingredients for a user
func GetAll(l *logrus.Entry, client *mongo.Client, userId string) ([]UserInventory, error) {
	collection := GetIngredientCollection(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userId": userId}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		l.WithError(err).Error("Failed to fetch user inventory")
		return nil, err
	}
	defer cursor.Close(ctx)

	// TODO Check if it's good
	inventories := make([]UserInventory, 0)
	if err = cursor.All(ctx, &inventories); err != nil {
		l.WithError(err).Error("Failed to decode user inventory")
		return nil, err
	}

	return inventories, nil
}

// GetOne retrieves a single inventory item by ID
func GetOne(l *logrus.Entry, client *mongo.Client, userId string, ingredientId string) (*UserInventory, error) {
	collection := GetIngredientCollection(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var inventory UserInventory
	err := collection.FindOne(ctx, bson.M{
		"ingredientId": ingredientId,
		"userId":       userId,
	}).Decode(&inventory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		l.WithError(err).Error("Failed to fetch inventory item")
		return nil, err
	}

	return &inventory, nil
}

// func FindByIdAndUnit(l *logrus.Entry, client *mongo.Client, id string, unit string) (*Ingredient, error) {
// 	coll := GetCollection(client, "inventory")
// 	filter := map[string]string{"ingredient_id": id, "units": unit}
// 	ingredient := Ingredient{}
// 	err := coll.FindOne(context.Background(), filter).Decode(&ingredient)
// 	if err != nil {
// 		l.WithError(err).Error("Error when trying to find ingredient by ID and unit")
// 		return nil, err
// 	}
// 	return &ingredient, nil
// }

// UpdateOne updates an existing inventory item
// TODO CHange the format of the update
func UpdateOne(l *logrus.Entry, client *mongo.Client, update *UserInventory) (*UserInventory, error) {
	collection := GetIngredientCollection(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	update.UpdatedAt = time.Now()

	filter := bson.M{
		"ingredientId": update.IngredientID,
		"userId":       update.UserID,
	}

	// Create a map to hold the fields to update
	updateFields := bson.M{
		"quantity":  update.Quantity,
		"unit":      update.Unit,
		"updatedAt": update.UpdatedAt,
	}

	// Only add name to update fields if it's not empty
	if update.Name != "" {
		updateFields["name"] = update.Name
	}

	updateDoc := bson.M{
		"$set": updateFields,
	}

	result := collection.FindOneAndUpdate(ctx, filter, updateDoc, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var updated UserInventory
	if err := result.Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		l.WithError(err).Error("Failed to update inventory item")
		return nil, err
	}

	return &updated, nil
}

// InsertOne creates a new inventory item
func InsertOne(l *logrus.Entry, client *mongo.Client, inventory *UserInventory) (*UserInventory, error) {
	collection := GetIngredientCollection(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	inventory.CreatedAt = now
	inventory.UpdatedAt = now

	result, err := collection.InsertOne(ctx, inventory)
	if err != nil {
		l.WithError(err).Error("Failed to insert inventory item")
		return nil, err
	}

	inventory.ID = result.InsertedID.(primitive.ObjectID)
	return inventory, nil
}

func DeleteOne(l *logrus.Entry, client *mongo.Client, userId string, ingredientId string) error {
	collection := GetIngredientCollection(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{
		"ingredientId": ingredientId,
		"userId":       userId,
	})
	if err != nil {
		l.WithError(err).Error("Failed to delete inventory item")
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
