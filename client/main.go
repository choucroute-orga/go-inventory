package main

import (
	"context"
	pb "inventory/inventory"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := "localhost:50051"
	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewInventoryClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateIngredient(ctx,
		&pb.PostIngredientRequest{
			Id:     "test-" + uuid.New().String(),
			Name:   "test",
			UserId: "gridexx",
			Amount: 1,
			Unit:   "i",
		})
	if err != nil {
		log.Fatalf("could not post ingredient: %v", err)
	}

	log.Printf("Ingredient added in inv: %s", r.GetId())

	r2, err := c.GetIngredient(ctx, &pb.GetIngredientRequest{
		IngredientId: r.GetId(),
		UserId:       "gridexx",
	})
	if err != nil {
		log.Fatalf("could not get ingredient: %v", err)
	}
	// Check that the response is what we expect.
	log.Printf("Ingredient: %s", r2.GetId())

	_, err = c.UpdateIngredient(ctx, &pb.PostIngredientRequest{
		Id: r.GetId(),

		Name:   "test",
		UserId: "gridexx",
		Amount: 2,
		Unit:   "g",
	})
	if err != nil {
		log.Fatalf("could not update ingredient: %v", err)
	}

	ri, err := c.GetIngredient(ctx, &pb.GetIngredientRequest{
		IngredientId: r.GetId(),
		UserId:       "gridexx",
	})

	if err != nil {
		log.Fatalf("could not get ingredient: %v", err)
	}

	// CHeck that the response is what we expect
	if ri.GetAmount() != 2 {
		log.Fatalf("Update: Amount should be 2, got %f", ri.GetAmount())
	}
	if ri.GetUnit() != "g" {
		log.Fatalf("Unit should be g, got %s", ri.GetUnit())
	}

	// Get the inventory to have the updated ingredient

	r5, err := c.GetUserInventory(ctx, &pb.GetInventoryRequest{
		UserId: "gridexx",
	})

	if err != nil {
		log.Fatalf("could not get inventory: %v", err)
	}

	if len(r5.UserInventory) != 1 {
		log.Printf("Inventory: %s", r5.String())
		log.Fatalf("Inventory should have 1 item, got %d", len(r5.UserInventory))
	}

	if r5.UserInventory[0].Amount != 2 {
		log.Fatalf("Amount should be 2, got %f", r5.UserInventory[0].Amount)
	}

	if r5.UserInventory[0].Unit != "g" {
		log.Fatalf("Unit should be g, got %s", r5.UserInventory[0].Unit)
	}

	log.Printf("Inventory: %s", r5.String())

	// Delete the ingredient
	r3, err := c.DeleteIngredient(ctx, &pb.DeleteIngredientRequest{
		IngredientId: r.GetId(),
		UserId:       "gridexx",
	})
	if err != nil {
		log.Fatalf("could not delete ingredient: %v", err)
	}

	log.Printf("Ingredient deleted: %s", r3.String())
}
