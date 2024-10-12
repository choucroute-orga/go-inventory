package main

import (
	"context"
	"fmt"
	"inventory/api"
	"inventory/configuration"
	"inventory/db"
	"inventory/messages"
	"inventory/validation"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var logger = logrus.WithFields(logrus.Fields{
	"context": "main",
})

func main() {
	configuration.SetupLogging()
	logger.Info("Inventory API Starting...")

	conf := configuration.New()
	logger.Logger.SetLevel(conf.LogLevel)
	mongo, err := db.New(conf)
	if err != nil {
		logger.WithError(err).Error("Unable to ping database to check connection.")
		return
	}

	coll := db.GetIngredientCollection(mongo)
	if coll == nil {
		logger.Error("Unable to get collection")
		return
	}

	val := validation.New(conf)
	r := api.New(val)
	v1 := r.Group(conf.ListenRoute)
	amqp := messages.New(conf)
	h := api.NewApiHandler(mongo, amqp, conf)

	h.Register(v1)
	tp, _ := api.InitOtel()
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.WithError(err).Error("Error shutting down tracer provider")
		}
		if err := mongo.Disconnect(context.TODO()); err != nil {
			logger.WithError(err).Error("Error closing mongo connection")
		}
		if err := amqp.Close(); err != nil {
			logger.WithError(err).Error("Error closing amqp connection")
		}
	}()

	go func() {
		h.ConsumeAddRecipeMessages(ctx)
	}()

	go func() {
		h.ConsumeAddIngredientMessage(ctx)
	}()

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		logger.Info("Shutting down gracefully...")
		cancel()
	}()

	if err := r.Start(fmt.Sprintf("%v:%v", conf.ListenAddress, conf.ListenPort)); err != nil {
		logger.WithError(err).Fatal("Error starting the server")
	}
}
