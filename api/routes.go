package api

import (
	"inventory/db"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var logger = logrus.WithField("context", "api/routes")

func (api *ApiHandler) getAliveStatus(c echo.Context) error {
	l := logger.WithField("request", "getAliveStatus")
	status := NewHealthResponse(LiveStatus)
	if err := c.Bind(status); err != nil {
		FailOnError(l, err)
		return NewInternalServerError(err.Error())
	}
	l.WithFields(logrus.Fields{
		"action": "getStatus",
		"status": status,
	}).Debug("Health Status ping")

	return c.JSON(http.StatusOK, &status)
}

func (api *ApiHandler) getReadyStatus(c echo.Context) error {
	l := logger.WithField("request", "getReadyStatus")
	err := api.dbh.Ping()
	WarnOnError(l, err)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, NewHealthResponse(NotReadyStatus))
	}

	return c.JSON(http.StatusOK, NewHealthResponse(ReadyStatus))
}

func (api *ApiHandler) getIngredients(c echo.Context) error {
	l := logger.WithField("request", "getIngredients")

	// TODO Check if we retrieve the userId from the request like this
	userId := c.QueryParam("userId")
	if userId == "" {
		return NewBadRequestError("userId is required")
	}

	inventories, err := api.dbh.GetUserInventory(l, userId)

	if err != nil {
		return NewInternalServerError(err.Error())
	}

	if inventories == nil {
		inventories = make([]db.UserInventory, 0)
	}

	return c.JSON(http.StatusOK, NewAllIngredientsResponse(&inventories))
}

// getIngredient handles GET /inventory/ingredient/:id
func (api *ApiHandler) getIngredient(c echo.Context) error {
	l := logger.WithField("request", "getIngredient")

	userId := c.QueryParam("userId")
	if userId == "" {
		return NewBadRequestError("userId is required")
	}

	id := c.Param("id")
	inventory, err := api.dbh.GetOneUserInventory(l, userId, id)
	if err != nil {
		return NewInternalServerError(err.Error())
	}

	if inventory == nil {
		return NewNotFoundError("Inventory item not found")
	}

	return c.JSON(http.StatusOK, NewIngredientResponse(inventory))
}

// insertOne handles POST /inventory/ingredient
func (api *ApiHandler) insertOne(c echo.Context) error {
	l := logger.WithField("request", "insertOne")

	var inventory PostIngredientRequest
	if err := c.Bind(&inventory); err != nil {
		return NewBadRequestError(err.Error())
	}

	if err := c.Validate(inventory); err != nil {
		return NewUnprocessableEntityError(err.Error())
	}

	inventoryIng, err := NewIngredientInventory(&inventory)
	if err != nil {
		return NewUnprocessableEntityError(err.Error())
	}
	result, err := api.dbh.InsertOneUserInventory(l, inventoryIng)
	if err != nil {
		return NewInternalServerError(err.Error())
	}

	return c.JSON(http.StatusCreated, NewIngredientResponse(result))
}

// updateOne handles PUT /inventory/ingredient/:id
func (api *ApiHandler) updateOne(c echo.Context) error {
	l := logger.WithField("request", "updateOne")

	var ingredient PutIngredientRequest
	if err := c.Bind(&ingredient); err != nil {
		return NewBadRequestError(err.Error())
	}

	if err := c.Validate(ingredient); err != nil {
		return NewUnprocessableEntityError(err.Error())
	}
	inventory, err := NewIngredientInventoryFromPut(&ingredient)
	if err != nil {
		return NewBadRequestError(err.Error())
	}
	result, err := api.dbh.UpdateOneUserInventory(l, inventory)
	if err != nil {
		return NewInternalServerError(err.Error())
	}

	if result == nil {
		return NewNotFoundError("Inventory item not found")
	}

	return c.JSON(http.StatusOK, NewIngredientResponse(result))
}

// deleteOne handles DELETE /inventory/ingredient/:id
func (api *ApiHandler) deleteOne(c echo.Context) error {
	l := logger.WithField("request", "deleteOne")
	var delete DeleteIngredientRequest
	if err := c.Bind(&delete); err != nil {
		return NewBadRequestError(err.Error())
	}
	if err := c.Validate(delete); err != nil {
		return NewUnprocessableEntityError(err.Error())
	}

	if err := api.dbh.DeleteOneUserInventory(l, delete.UserID, delete.ID); err != nil {
		if err == mongo.ErrNoDocuments {
			return NewNotFoundError("Inventory item not found")
		}
		return NewInternalServerError(err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
