package api

import (
	"fmt"
	"inventory/db"
	"inventory/messages"
	"inventory/tests"
	"log"
	"math/rand"
	"testing"

	"github.com/sirupsen/logrus"
)

func doubleMe(x float64) float64 {
	return x * 2
}

// You can use testing.T, if you want to test the code without benchmarking
func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup suite")

	// Return a function to teardown the test
	return func(tb testing.TB) {
		log.Println("teardown suite")
	}
}

// Almost the same as the above, but this one is for single test instead of collection of tests
func setupTest(tb testing.TB) (*ApiHandler, func(tb testing.TB)) {
	// log.Println("setup test")

	// return func(tb testing.TB) {
	// 	log.Println("teardown test")
	// }

	// Get a random port for the test, between 1024 and 65535
	exposedPort := fmt.Sprint(rand.Intn(65525-1024) + 1024)
	mongo, pool, resource := tests.InitTestDocker(exposedPort)
	conf := tests.GetDefaultConf()
	api := NewApiHandler(mongo, nil, conf)
	tests.SeedDatabase(mongo)
	return api, func(tb testing.TB) {
		tests.CloseTestDocker(mongo, pool, resource)
	}
}

func TestDB(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	t.Run("Insert ingredient in the DB", func(t *testing.T) {
		logrus.SetLevel(logrus.DebugLevel)
		l := logrus.WithField("test", "Insert ingredient in the DB")
		api, teardownTest := setupTest(t)

		id := "665d9ea9d0ba48dcb11ac3a1"

		i1 := db.Ingredient{
			ID:           db.NewID(),
			IngredientID: id,
			Name:         "Ingredient 1",
			Quantity:     20,
			Units:        "g",
		}

		i2 := db.Ingredient{
			ID:           db.NewID(),
			IngredientID: "665d9ea9d0ba48dcb11ac3a1",
			Name:         "Ingredient 1",
			Quantity:     1,
			Units:        "kg",
		}

		db.InsertOne(l, api.mongo, i1)
		db.InsertOne(l, api.mongo, i2)

		ingredients, err := db.FindById(l, api.mongo, id)

		// Err should be nil
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if len(*ingredients) != 2 {
			t.Errorf("Expected 2 ingredients, got %d", len(*ingredients))
		}

		// Test the update
		i1.Quantity = 40
		db.UpdateOne(l, api.mongo, i1)
		i, err := db.FindByIdAndUnit(l, api.mongo, id, "g")

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		// Ensure the update was successful
		if i.Quantity != 40 {
			t.Errorf("Expected 40, got %f", i.Quantity)
		}
		if i.Units != "g" {
			t.Errorf("Expected g, got %s", i.Units)
		}
		if i.IngredientID != id {
			t.Errorf("Expected %s, got %s", id, i.IngredientID)
		}
		if i.Name != "Ingredient 1" {
			t.Errorf("Expected Ingredient 1, got %s", i.Name)
		}
		defer teardownTest(t)
	})

	t.Run("Insert ingredient and delete in the DB", func(t *testing.T) {
		logrus.SetLevel(logrus.DebugLevel)
		l := logrus.WithField("test", "Insert ingredient and delete in the DB")
		api, teardownTest := setupTest(t)
		i := db.Ingredient{
			ID:           db.NewID(),
			IngredientID: "665d9ea9d0ba48dcb11ac3a1",
			Name:         "Ingredient 1",
			Quantity:     20,
			Units:        "g",
		}
		err := db.InsertOne(l, api.mongo, i)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		err = db.DeleteOne(l, api.mongo, i.IngredientID)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		defer teardownTest(t)
	})

	t.Run("Insert one Recipe in the DB", func(t *testing.T) {
		logrus.SetLevel(logrus.DebugLevel)
		l := logrus.WithField("test", "Insert one Recipe in the DB")
		api, teardownTest := setupTest(t)

		addRecipe := messages.AddRecipe{
			ID: "1",
			Ingredients: []messages.AddIngredient{
				{
					ID:     "665d9ea9d0ba48dcb11ac3a1",
					Amount: 100,
					Unit:   "g",
				},
				{
					ID:     "665d9ea9d0ba48dcb11ac3a2",
					Amount: 200,
					Unit:   "g",
				},
				{
					ID:     "665d9ea9d0ba48dcb11ac3a3",
					Amount: 3,
					Unit:   "kg",
				},
				{
					ID:     "665d9ea9d0ba48dcb11ac3a4",
					Amount: 2,
					Unit:   "kg",
				},
			},
		}

		id1 := "665d9ea9d0ba48dcb11ac3a1"
		id2 := "665d9ea9d0ba48dcb11ac3a2"
		id3 := "665d9ea9d0ba48dcb11ac3a3"

		r1 := db.Ingredient{
			ID:           db.NewID(),
			IngredientID: id1,
			Name:         "Ingredient 1",
			Quantity:     20,
			Units:        "g",
		}

		r2 := db.Ingredient{
			ID:           db.NewID(),
			IngredientID: id2,
			Name:         "Ingredient 2",
			Quantity:     20,
			Units:        "kg",
		}

		r3 := db.Ingredient{
			ID:           db.NewID(),
			IngredientID: id3,
			Name:         "Ingredient 3",
			Quantity:     20,
			Units:        "kg",
		}

		db.InsertOne(l, api.mongo, r1)
		db.InsertOne(l, api.mongo, r2)
		db.InsertOne(l, api.mongo, r3)
		iSl := *api.createShoppingList(l, addRecipe)

		// Check if the shopping list is correct
		if len(iSl) != 3 {
			t.Errorf("Expected 3 ingredients in the shopping list, got %d", len(iSl))
		}
		if iSl[0].Amount != 80 {
			t.Errorf("Expected amount of 80, got %f", iSl[0].Amount)
		}
		if iSl[1].Amount != 200 {
			t.Errorf("Expected amount of 200, got %f", iSl[1].Amount)
		}
		if iSl[2].Amount != 2 {
			t.Errorf("Expected amount of 2, got %f", iSl[2].Amount)
		}
		defer teardownTest(t)
	})

}
