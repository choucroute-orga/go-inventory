package api

import (
	"inventory/configuration"
	"inventory/db"
	"inventory/validation"

	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ApiHandler struct {
	dbh        db.DbHandler
	conf       *configuration.Configuration
	amqp       *amqp.Connection
	validation *validation.Validation
	tracer     trace.Tracer
}

func NewApiHandler(dbh db.DbHandler, amqp *amqp.Connection, conf *configuration.Configuration) *ApiHandler {
	handler := ApiHandler{
		dbh:        dbh,
		amqp:       amqp,
		conf:       conf,
		validation: validation.New(conf),
		tracer:     otel.Tracer(conf.OtelServiceName),
	}
	return &handler
}

func (api *ApiHandler) Register(v1 *echo.Group) {

	health := v1.Group("/health")
	health.GET("/alive", api.getAliveStatus)
	health.GET("/live", api.getAliveStatus)
	health.GET("/ready", api.getReadyStatus)

	inventory := v1.Group("/inventory")
	inventory.GET("/ingredient", api.getIngredients)
	inventory.GET("/ingredient/:id", api.getIngredient)
	inventory.POST("/ingredient", api.insertOne)
	inventory.PUT("/ingredient/:id", api.updateOne)
	inventory.DELETE("/ingredient/:id/user/:userId", api.deleteOne)
}
