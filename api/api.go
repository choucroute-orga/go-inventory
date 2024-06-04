package api

import (
	"inventory/configuration"
	"inventory/validation"

	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiHandler struct {
	mongo      *mongo.Client
	conf       *configuration.Configuration
	amqp       *amqp.Connection
	validation *validation.Validation
}

func NewApiHandler(mongo *mongo.Client, amqp *amqp.Connection, conf *configuration.Configuration) *ApiHandler {
	handler := ApiHandler{
		mongo:      mongo,
		amqp:       amqp,
		conf:       conf,
		validation: validation.New(conf),
	}
	return &handler
}

func (api *ApiHandler) Register(v1 *echo.Group) {

	health := v1.Group("/health")
	health.GET("/alive", api.getAliveStatus)
	health.GET("/live", api.getAliveStatus)
	health.GET("/ready", api.getReadyStatus)

	inventory := v1.Group("/inventory")
	inventory.GET("/ingredient/:id", api.getIngredient)
	inventory.POST("/ingredient", api.insertOne)

}
