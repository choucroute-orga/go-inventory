package main

import (
	"context"
	"fmt"
	"inventory/api"
	"inventory/configuration"
	"inventory/db"
	"inventory/validation"

	"github.com/sirupsen/logrus"
)

var logger = logrus.WithFields(logrus.Fields{
	"context": "main",
})

func main() {
	logger.Info("Cacahuete API Starting...")

	conf := configuration.New()

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

	// Insert an ingredient
	// ingredient := db.Ingredient{
	// 	Name: "Cacahuete",
	// 	Unit: "g",
	// }
	// _, err = coll.InsertOne(context.Background(), ingredient)
	// if err != nil {
	// 	logger.WithError(err).Error("Unable to insert ingredient")
	// 	return
	// }

	val := validation.New(conf)
	r := api.New(val)
	v1 := r.Group(conf.ListenRoute)

	h := api.NewApiHandler(mongo, conf)

	h.Register(v1, conf)
	r.Logger.Fatal(r.Start(fmt.Sprintf("%v:%v", conf.ListenAddress, conf.ListenPort)))

	defer func() {
		if err := mongo.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}
