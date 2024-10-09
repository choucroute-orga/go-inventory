package main

import (
	"context"
	"fmt"
	"inventory/api"
	"inventory/configuration"
	"inventory/db"
	"inventory/messages"
	"inventory/validation"

	"github.com/sirupsen/logrus"
)

var logger = logrus.WithFields(logrus.Fields{
	"context": "main",
})

func main() {
	logger.Info("Cacahuete API Starting...")

	conf := configuration.New()
	logger.Logger.SetLevel(conf.LogLevel)
	mongo, err := db.New(conf)
	if err != nil {
		logger.WithError(err).Error("Unable to ping database to check connection.")
		return
	}

	coll := db.GetCollection(mongo, "inventory")
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
	go func() {
		r.Logger.Fatal(r.Start(fmt.Sprintf("%v:%v", conf.ListenAddress, conf.ListenPort)))
	}()

	h.ConsumesMessages()

	defer func() {
		if err := mongo.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		if err := amqp.Close(); err != nil {
			panic(err)
		}
	}()
}
