package api

import (
	"context"
	"errors"
	"inventory/db"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var logger = logrus.WithField("context", "api/routes")

func (api *ApiHandler) getAliveStatus(c echo.Context) error {
	l := logger.WithField("request", "getAliveStatus")
	status := NewHealthResponse(LiveStatus)
	if err := c.Bind(status); err != nil {
		FailOnError(l, err, "Response binding failed")
		return NewInternalServerError(err)
	}
	l.WithFields(logrus.Fields{
		"action": "getStatus",
		"status": status,
	}).Debug("Health Status ping")

	return c.JSON(http.StatusOK, &status)
}

func (api *ApiHandler) getReadyStatus(c echo.Context) error {
	l := logger.WithField("request", "getReadyStatus")
	err := db.GetCollection(api.mongo, "inventory").Database().Client().Ping(context.Background(), nil)
	if err != nil {
		WarnOnError(l, err, "Unable to ping database to check connection.")
		return c.JSON(http.StatusServiceUnavailable, NewHealthResponse(NotReadyStatus))
	}

	return c.JSON(http.StatusOK, NewHealthResponse(ReadyStatus))
}

func (api *ApiHandler) getIngredient(c echo.Context) error {
	l := logger.WithField("request", "getIngredient")

	// Retrieve the name from the request
	id := c.Param("id")
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewBadRequestError(err)
	}
	// Retrieve the ingredient from the database
	ingredient, err := db.FindById(l, api.mongo, idObject)
	if err != nil {
		return NewNotFoundError(errors.New("ingredient not found"))
	}
	return c.JSON(http.StatusOK, ingredient)
}

func (api *ApiHandler) insertOne(c echo.Context) error {
	l := logger.WithField("request", "insertOne")

	// Retrieve the name from the request

	var ingredient IngredientRequest
	if err := c.Bind(&ingredient); err != nil {
		return NewBadRequestError(err)
	}
	if err := c.Validate(ingredient); err != nil {
		return NewBadRequestError(err)
	}

	id, err := primitive.ObjectIDFromHex(ingredient.ID)

	if err != nil {
		return NewBadRequestError(err)
	}

	// TODO Add the quantity and units to the ingredient in an array
	ingredientDB := db.Ingredient{
		ID:       id,
		Name:     ingredient.Name,
		Quantity: ingredient.Quantity,
		Units:    ingredient.Units,
	}

	// Insert the ingredient into the database
	err = db.InsertOne(l, api.mongo, ingredientDB)
	if err != nil {
		return NewInternalServerError(err)
	}
	return c.JSON(http.StatusOK, ingredient)

}
