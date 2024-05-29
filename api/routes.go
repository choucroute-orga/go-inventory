package api

import (
	"context"
	"errors"
	"inventory/db"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
	name := c.Param("name")
	// Retrieve the ingredient from the database
	ingredient := db.FindOne(l, api.mongo, name)
	if ingredient.Name == "" {
		return NewNotFoundError(errors.New("ingredient not found"))
	}
	return c.JSON(http.StatusOK, ingredient)
}

func (api *ApiHandler) insertOne(c echo.Context) error {
	l := logger.WithField("request", "insertOne")

	// Retrieve the name from the request
	name := c.Param("name")
	ingredient := db.FindOne(l, api.mongo, name)
	if ingredient.Name != "" {
		return NewConflictError(errors.New("ingredient already exists"))
	}
	// Define a random ID for the ingredient
	id := db.NewID()
	ingredient = db.Ingredient{
		ID:   id,
		Name: name,
		Unit: "g",
	}

	// Insert the ingredient into the database
	err := db.InsertOne(l, api.mongo, ingredient)
	if err != nil {
		return NewInternalServerError(err)
	}
	return c.JSON(http.StatusOK, ingredient)

}
