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

var loger = logrus.WithFields(logrus.Fields{
	"context": "db/query",
})

func NewID() primitive.ObjectID {
	return primitive.NewObjectIDFromTimestamp(time.Now())
}

func GetAll(l *logrus.Entry, client *mongo.Client) ([]Ingredient, error) {
	coll := GetCollection(client, "inventory")
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		l.WithError(err).Error("Error when trying to find all ingredients")
		return nil, err
	}
	ingredients := make([]Ingredient, 0)
	err = cursor.All(context.Background(), &ingredients)
	if err != nil {
		l.WithError(err).Error("Error when trying to decode all ingredients")
		return nil, err
	}
	return ingredients, nil
}

func FindById(l *logrus.Entry, client *mongo.Client, id string) (*[]Ingredient, error) {
	coll := GetCollection(client, "inventory")
	ingredients := make([]Ingredient, 0)
	filter := map[string]string{"ingredient_id": id}
	cursor, err := coll.Find(context.Background(), filter)

	if err != nil {
		l.WithError(err).Error("Error when trying to find ingredient by ID")
		return nil, err
	}
	// Decode all the ingredients
	err = cursor.All(context.Background(), &ingredients)
	if err != nil {
		l.WithError(err).Error("Error when trying to decode all ingredients")
		return nil, err
	}
	return &ingredients, nil
}

func FindByIdAndUnit(l *logrus.Entry, client *mongo.Client, id string, unit string) (*Ingredient, error) {
	coll := GetCollection(client, "inventory")
	filter := map[string]string{"ingredient_id": id, "units": unit}
	ingredient := Ingredient{}
	err := coll.FindOne(context.Background(), filter).Decode(&ingredient)
	if err != nil {
		l.WithError(err).Error("Error when trying to find ingredient by ID and unit")
		return nil, err
	}
	return &ingredient, nil
}

func UpdateOne(l *logrus.Entry, client *mongo.Client, ingredient Ingredient) error {
	coll := GetCollection(client, "inventory")
	filter := map[string]string{"ingredient_id": ingredient.IngredientID}
	update := bson.M{"$set": ingredient}
	res, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		l.WithError(err).Error("Error when trying to update ingredient")
		return err
	}
	if res.MatchedCount == 0 {
		err = errors.New("ID not found")
		l.WithError(err).Error("Error when trying to update ingredient")
		return err
	}
	return nil
}

func InsertOne(l *logrus.Entry, client *mongo.Client, ingredient Ingredient) error {
	coll := GetCollection(client, "inventory")
	_, err := coll.InsertOne(context.TODO(), ingredient)
	if err != nil {
		l.WithError(err).Error("Error when trying to insert the " + ingredient.Name + " ingredient")
		return err
	}
	return nil
}

// Delete with the ID instead of IngredientID
func DeleteOne(l *logrus.Entry, client *mongo.Client, id string) error {
	coll := GetCollection(client, "inventory")
	filter := map[string]string{"ingredient_id": id}
	row, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		l.WithError(err).Error("Error when trying to delete the " + id + " ingredient")
		return err
	}
	if row.DeletedCount == 0 {
		err := errors.New("ingredient id not found")
		l.Error(err.Error())
		return err
	}
	return nil
}
