package api

import (
	"context"
	"fmt"
	"inventory/db"
	"inventory/messages"
	"inventory/tests"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	dockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient   *mongo.Client
	mongoPool     *dockertest.Pool
	mongoResource *dockertest.Resource
	once          sync.Once
)

// InitTestMongo initializes a single MongoDB instance for all tests
func InitTestMongo() (*mongo.Client, error) {
	var initErr error
	once.Do(func() {
		// Create a new pool
		pool, err := dockertest.NewPool("")
		if err != nil {
			initErr = fmt.Errorf("could not construct pool: %w", err)
			return
		}

		mongoPool = pool

		// Set a timeout for docker operations
		pool.MaxWait = time.Second * 30

		// Start MongoDB container
		resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "mongo",
			Tag:        "5.0",
			Env: []string{
				"MONGO_INITDB_ROOT_USERNAME=root",
				"MONGO_INITDB_ROOT_PASSWORD=password",
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})

		if err != nil {
			initErr = fmt.Errorf("could not start resource: %w", err)
			return
		}

		mongoResource = resource

		// Initialize mongo client
		mongoURI := fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp"))

		// Retry connection with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				initErr = fmt.Errorf("timeout waiting for mongodb to be ready")
				return
			case <-ticker.C:
				client, err := mongo.Connect(
					context.Background(),
					options.Client().ApplyURI(mongoURI).SetConnectTimeout(2*time.Second),
				)
				if err != nil {
					continue
				}

				// Try to ping
				if err := client.Ping(context.Background(), nil); err != nil {
					_ = client.Disconnect(context.Background())
					continue
				}

				mongoClient = client
				return
			}
		}
	})

	return mongoClient, initErr
}

// CleanupDatabase removes all data from the test database
func CleanupDatabase(t *testing.T, client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("inventory").Collection("inventory").DeleteMany(ctx, bson.M{})

	if err != nil {
		t.Logf("Warning: Failed to cleanup database: %v", err)
	}

	err = client.Database("inventory").Drop(ctx)

	// Remove the inventory collection

	if err != nil {
		t.Logf("Warning: Failed to cleanup database: %v", err)
	}
}

func setupTest(t *testing.T) (*ApiHandler, func()) {
	t.Helper()

	// Use existing MongoDB instance
	client := mongoClient
	if client == nil {
		t.Fatal("MongoDB client not initialized")
	}

	// Clean the database
	CleanupDatabase(t, client)

	// Initialize the database
	tests.SeedDatabase(client)

	// Create API handler
	conf := tests.GetDefaultConf()
	dbh := db.NewMongoHandler(client)
	api := NewApiHandler(dbh, nil, conf)

	// Return cleanup function
	return api, func() {
		CleanupDatabase(t, client)
	}
}

func TestMain(m *testing.M) {
	// Setup
	client, err := InitTestMongo()
	if err != nil {
		log.Fatalf("Could not start MongoDB: %s", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(ctx)
	}
	if mongoPool != nil && mongoResource != nil {
		_ = mongoPool.Purge(mongoResource)
	}

	os.Exit(code)
}

func TestDB(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Insert ingredient in the DB",
			test: func(t *testing.T) {
				api, cleanup := setupTest(t)
				l := logrus.WithField("test", "Insert ingredient in the DB")
				defer cleanup()
				id := "605d9ea9d0ba48dcb11ac3a1"

				i1 := db.UserInventory{
					ID:           db.NewID(),
					UserID:       "1",
					IngredientID: id,
					Name:         "UserInventory 1",
					Quantity:     20,
					Unit:         "g",
				}
				api.dbh.InsertOneUserInventory(l, &i1)
				ingredient, err := api.dbh.GetOneUserInventory(l, i1.UserID, i1.IngredientID)

				// Err should be nil
				if err != nil {
					t.Errorf("Expected nil, got %v", err)
				}

				// Check if the ingredient is correct
				if ingredient.Quantity != 20 {
					t.Errorf("Expected 20, got %f", ingredient.Quantity)
				}

				// Test the update
				i1.Quantity = 40
				i1.Name = "UserInventory 1 Modified"
				api.dbh.UpdateOneUserInventory(l, &i1)
				i, err := api.dbh.GetOneUserInventory(l, i1.UserID, i1.IngredientID)

				if err != nil {
					t.Errorf("Expected nil, got %v", err)
				}
				// Ensure the update was successful
				if i.Quantity != 40 {
					t.Errorf("Expected 40, got %f", i.Quantity)
				}
				if i.Unit != "g" {
					t.Errorf("Expected g, got %s", i.Unit)
				}
				if i.IngredientID != id {
					t.Errorf("Expected %s, got %s", id, i.IngredientID)
				}
				if i.Name != "UserInventory 1 Modified" {
					t.Errorf("Expected UserInventory 1, got %s", i.Name)
				}
			},
		},
		{
			name: "Insert ingredient and delete in the DB",
			test: func(t *testing.T) {
				api, cleanup := setupTest(t)
				l := logrus.WithField("test", "Insert ingredient and delete in the DB")
				defer cleanup()

				i := db.UserInventory{
					ID:           db.NewID(),
					UserID:       "1",
					IngredientID: "a65d9ea9d0ba48dcb11ac3a1",
					Name:         "UserInventory 1",
					Quantity:     20,
					Unit:         "g",
				}
				_, err := api.dbh.InsertOneUserInventory(l, &i)
				if err != nil {
					t.Errorf("Expected nil, got %v", err)
				}
				err = api.dbh.DeleteOneUserInventory(l, i.UserID, i.IngredientID)
				if err != nil {
					t.Errorf("Expected nil, got %v", err)
				}
				_, err = api.dbh.GetOneUserInventory(l, i.UserID, i.IngredientID)
				if err != mongo.ErrNoDocuments {
					t.Errorf("Expected mongo.ErrNoDocuments, got %v", err)
				}
			},
		},
		{
			name: "Check conversion to base unit",
			test: func(t *testing.T) {
				_, cleanup := setupTest(t)
				defer cleanup()

				res, err := ConvertToBaseUnit(0.5, "kg")
				if err != nil {
					t.Errorf("Expected nil, got %v", err)
				}
				if res.Quantity != 500 && res.Unit != "g" {
					t.Errorf("Expected 500g, got %f%s", res.Quantity, res.Unit)
				}

				res, err = ConvertToBaseUnit(0.5, "l")
				if err != nil {
					t.Errorf("Expected nil, got %v", err)
				}
				if res.Quantity != 500 && res.Unit != "ml" {
					t.Errorf("Expected 500ml, got %f%s", res.Quantity, res.Unit)
				}
				res, err = ConvertToBaseUnit(3, "is")
				if err != nil {
					t.Errorf("Expected nil, got %v", err)
				}
				if res.Quantity != 3 && res.Unit != "i" {
					t.Errorf("Expected 3i, got %f%s", res.Quantity, res.Unit)
				}

				// Check conversion with invalid unit
				_, err = ConvertToBaseUnit(0.5, "invalid")
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

				// Check Ratio Conversion with cs, tbsp, tsp
				res, err = ConvertToBaseUnit(1, "cup")
				if err != nil {
					t.Errorf("Expected nil, got %v", err)
				}
				if res.Quantity != 236.588 && res.Unit != "ml" {
					t.Errorf("Expected 236.588ml, got %f%s", res.Quantity, res.Unit)
				}
			},
		},
		{
			name: "Create the Shopping List",
			test: func(t *testing.T) {
				l := logrus.WithField("test", "Insert one Recipe in the DB")
				api, cleanup := setupTest(t)
				defer cleanup()

				id1 := "325d9ea9d0ba48dcb11ac3a1"
				id2 := "325d9ea9d0ba48dcb11ac3a2"
				id3 := "325d9ea9d0ba48dcb11ac3a3"
				userId := "2xezfZ"
				addRecipe := messages.AddRecipe{
					ID:     "1",
					UserID: userId,
					Ingredients: []messages.Ingredient{
						{
							ID:     id1,
							Amount: 100,
							Unit:   "g",
						},
						{
							ID:     id2,
							Amount: 1,
							Unit:   "kg",
						},
						{
							ID:     id3,
							Amount: 2000,
							Unit:   "g",
						},
						{
							ID:     "665d9eh7iEba48dcb11ac3a4",
							Amount: 2,
							Unit:   "kg",
						},
					},
				}

				r1 := &db.UserInventory{
					ID:           db.NewID(),
					UserID:       userId,
					IngredientID: id1,
					Name:         "UserInventory 1",
					Quantity:     20,
					Unit:         "g",
				}

				r2 := &db.UserInventory{
					ID:           db.NewID(),
					UserID:       userId,
					IngredientID: id2,
					Name:         "UserInventory 2",
					Quantity:     500,
					Unit:         "g",
				}

				r3 := &db.UserInventory{
					ID:           db.NewID(),
					UserID:       userId,
					IngredientID: id3,
					Name:         "UserInventory 3",
					Quantity:     0.70,
					Unit:         "kg",
				}

				api.dbh.InsertOneUserInventory(l, r1)
				api.dbh.InsertOneUserInventory(l, r2)
				api.dbh.InsertOneUserInventory(l, r3)
				ctx := context.Background()
				iSl := *api.createShoppingList(ctx, l, addRecipe)

				// Check if the shopping list is correct
				if len(iSl) != 4 {
					t.Errorf("Expected 4 ingredients in the shopping list, got %d", len(iSl))
				}
				if iSl[0].Amount != 80 {
					t.Errorf("Expected amount of 80, got %f", iSl[0].Amount)
				}
				if iSl[1].Amount != 500 {
					t.Errorf("Expected amount of 500, got %f", iSl[1].Amount)
				}
				if iSl[2].Amount != 1.3 {
					t.Errorf("Expected amount of 1.3, got %f", iSl[2].Amount)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.test(t)
		})
	}
}
