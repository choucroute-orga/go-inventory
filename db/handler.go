package db

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DbHandler interface {
	Disconnect() error
	Ping() error
	NewID() primitive.ObjectID
	GetUserInventory(l *logrus.Entry, userId string) ([]UserInventory, error)
	GetOneUserInventory(l *logrus.Entry, userId string, ingredientId string) (*UserInventory, error)
	InsertOneUserInventory(l *logrus.Entry, inventory *UserInventory) (*UserInventory, error)
	UpdateOneUserInventory(l *logrus.Entry, udpate *UserInventory) (*UserInventory, error)
	DeleteOneUserInventory(l *logrus.Entry, userId string, ingredientId string) error
}

type MongoHandler struct {
	client                   *mongo.Client
	ingredientCollectionName string
	database                 string
}

func NewMongoHandler(client *mongo.Client) *MongoHandler {
	return &MongoHandler{
		client:                   client,
		ingredientCollectionName: "ingredient",
		database:                 "inventory",
	}
}

func (m *MongoHandler) getIngredientCollection() *mongo.Collection {
	return m.client.Database(m.database).Collection(m.ingredientCollectionName)
}

func (m *MongoHandler) Disconnect() error {
	return m.client.Disconnect(context.Background())
}

func (m *MongoHandler) Ping() error {
	coll := m.getIngredientCollection()
	if coll == nil {
		return errors.New("Failed to get collection")
	}
	return coll.Database().Client().Ping(context.Background(), nil)
}

func (m *MongoHandler) NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (m *MongoHandler) GetUserInventory(l *logrus.Entry, userId string) ([]UserInventory, error) {
	collection := m.getIngredientCollection()
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

func (m *MongoHandler) GetOneUserInventory(l *logrus.Entry, userId string, ingredientId string) (*UserInventory, error) {
	collection := m.getIngredientCollection()
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

func (m *MongoHandler) InsertOneUserInventory(l *logrus.Entry, inventory *UserInventory) (*UserInventory, error) {
	collection := m.getIngredientCollection()
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

func (m *MongoHandler) UpdateOneUserInventory(l *logrus.Entry, update *UserInventory) (*UserInventory, error) {
	collection := m.getIngredientCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"ingredientId": update.IngredientID,
		"userId":       update.UserID,
	}

	res, err := collection.ReplaceOne(ctx, filter, update)
	if err != nil {
		l.WithError(err).Error("Failed to update inventory item")
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return update, nil
}

func (m *MongoHandler) DeleteOneUserInventory(l *logrus.Entry, userId string, ingredientId string) error {
	collection := m.getIngredientCollection()
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
