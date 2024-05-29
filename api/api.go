package api

import (
	"inventory/configuration"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiHandler struct {
	mongo *mongo.Client
	conf  *configuration.Configuration
}

func NewApiHandler(mongo *mongo.Client, conf *configuration.Configuration) *ApiHandler {
	handler := ApiHandler{
		mongo: mongo,
		conf:  conf,
	}
	return &handler
}

func (api *ApiHandler) Register(v1 *echo.Group, conf *configuration.Configuration) {

	health := v1.Group("/health")
	health.GET("/alive", api.getAliveStatus)
	health.GET("/live", api.getAliveStatus)
	health.GET("/ready", api.getReadyStatus)

	inventory := v1.Group("/inventory")
	inventory.GET("/ingredient/:name", api.getIngredient)
	inventory.POST("/ingredient/:name", api.insertOne)

}
