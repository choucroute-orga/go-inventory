package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var loger = logrus.WithFields(logrus.Fields{
	"context": "db/query",
})

// func LogAndReturnError(l *logrus.Entry, result *gorm.DB, action string, modelType string) error {
// 	if err := result.Error; err != nil {
// 		l.WithError(err).Error("Error when trying to query database to " + action + " " + modelType)
// 		return err
// 	}
// 	return nil
// }

func NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func FindById(l *logrus.Entry, client *mongo.Client, id primitive.ObjectID) (Ingredient, error) {
	coll := GetCollection(client, "inventory")
	var ingredient Ingredient
	err := coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&ingredient)
	if err != nil {
		l.WithError(err).Error("Error when trying to find the " + id.Hex() + " ingredient")
	}
	return ingredient, err
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
