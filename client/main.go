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
		UserId: 	 "gridexx",
	})
	if err != nil {
		log.Fatalf("could not get ingredient: %v", err)
	}
	// Check that the response is what we expect.
	log.Printf("Ingredient: %s", r2.GetId())

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
