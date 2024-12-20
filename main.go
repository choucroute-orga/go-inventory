package main

import (
	"context"
	"fmt"
	"inventory/api"
	"inventory/configuration"
	"inventory/db"
	"inventory/server"

	// "inventory/grpc"
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
	dbh := db.NewMongoHandler(mongo)
	if err := dbh.Ping(); err != nil {
		logger.WithError(err).Error("Unable to ping database to check connection.")
		return
	}
	val := validation.New(conf)
	r := api.New(val)
	v1 := r.Group(conf.ListenRoute)
	amqp := messages.New(conf)
	h := api.NewApiHandler(dbh, amqp, conf)

	h.Register(v1)
	tp := api.InitOtel()
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

	go func() {
		logger.Info("Started gRPC server on port ", conf.GrpcPort)
		server.Run(dbh, conf.GrpcPort, conf)
	}()

	if err := r.Start(fmt.Sprintf("%v:%v", conf.ListenAddress, conf.ListenPort)); err != nil {
		logger.WithError(err).Fatal("Error starting the server")
	}
}
