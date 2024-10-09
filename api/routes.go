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

func (api *ApiHandler) getIngredients(c echo.Context) error {
	l := logger.WithField("request", "getIngredients")
	ingredients, err := db.GetAll(l, api.mongo)
	if err != nil {
		return NewInternalServerError(err)
	}
	return c.JSON(http.StatusOK, ingredients)
}

func (api *ApiHandler) getIngredient(c echo.Context) error {
	l := logger.WithField("request", "getIngredient")

	// Retrieve the name from the request
	id := c.Param("id")

	// Retrieve the ingredient from the database
	ingredients, err := db.FindById(l, api.mongo, id)
	if err != nil {
		return NewNotFoundError(errors.New("ingredient not found"))
	}
	ingredient := NewIngredientResponse(ingredients)
	return c.JSON(http.StatusOK, &ingredient)
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

	ingredientDB, _ := db.FindByIdAndUnit(l, api.mongo, ingredient.ID, ingredient.Unit)
	if ingredientDB != nil {
		ingredientDB.Quantity += ingredient.Amount
		if err := db.UpdateOne(l, api.mongo, *ingredientDB); err != nil {
			return NewInternalServerError(err)
		}
		return c.JSON(http.StatusOK, ingredient)
	}

	// TODO Add the quantity and units to the ingredient in an array
	ingredientDB = &db.Ingredient{
		ID:           db.NewID(),
		IngredientID: ingredient.ID,
		Name:         ingredient.Name,
		Quantity:     ingredient.Amount,
		Units:        ingredient.Unit,
	}

	// Insert the ingredient into the database
	if err := db.InsertOne(l, api.mongo, *ingredientDB); err != nil {
		return NewInternalServerError(err)
	}
	return c.JSON(http.StatusOK, ingredient)

}

// Delete by quantity and unit
func (api *ApiHandler) deleteOne(c echo.Context) error {
	l := logger.WithField("request", "deleteOne")

	// Retrieve the name from the request
	id := c.Param("id")

	// Delete the ingredient from the database
	if err := db.DeleteOne(l, api.mongo, id); err != nil {
		return NewInternalServerError(err)
	}
	return c.NoContent(http.StatusNoContent)
}
