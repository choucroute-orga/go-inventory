package server

import (
	"context"
	"errors"
	"fmt"
	"inventory/api"
	"inventory/configuration"
	"inventory/db"
	pb "inventory/inventory"
	"inventory/messages"
	"inventory/validation"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var logger = logrus.WithField("context", "grpc/server")

type server struct {
	pb.UnimplementedInventoryServer
	dbh        db.DbHandler
	validation *validation.Validation
}

func NewIngredientResponse(ingredient *db.UserInventory) *pb.IngredientResponse {
	return &pb.IngredientResponse{
		Id:        ingredient.IngredientID,
		Name:      ingredient.Name,
		Amount:    ingredient.Quantity,
		Unit:      ingredient.Unit.String(),
		CreatedAt: timestamppb.New(ingredient.CreatedAt),
		UpdatedAt: timestamppb.New(ingredient.UpdatedAt),
	}
}

func NewAllIngredientsResponse(ingredients *[]db.UserInventory) *pb.GetUserInventoryResponse {
	ingredientsResponse := []*pb.IngredientResponse{}

	for _, ingredient := range *ingredients {
		ingredientsResponse = append(ingredientsResponse, NewIngredientResponse(&ingredient))
	}
	if ingredientsResponse == nil {
		ingredientsResponse = make([]*pb.IngredientResponse, 0)
	}
	return &pb.GetUserInventoryResponse{
		UserInventory: ingredientsResponse,
	}
}

func (s *server) GetUserInventory(ctx context.Context, in *pb.GetInventoryRequest) (*pb.GetUserInventoryResponse, error) {
	inventory, err := s.dbh.GetUserInventory(logger, in.GetUserId())

	if err != nil {
		logger.WithError(err).Error("Error getting inventory")
		return nil, err
	}
	return NewAllIngredientsResponse(&inventory), nil

}

func (s *server) GetIngredient(ctx context.Context, in *pb.GetIngredientRequest) (*pb.IngredientResponse, error) {
	ingredient, err := s.dbh.GetOneUserInventory(logger, in.GetUserId(), in.GetIngredientId())
	if err != nil {
		logger.WithError(err).Error("Error getting inventory")
		return nil, err
	}

	if ingredient == nil {
		return nil, errors.New("No ingredient found")
	}

	return NewIngredientResponse(ingredient), nil
}

func (s *server) CreateIngredient(ctx context.Context, in *pb.PostIngredientRequest) (*pb.IngredientResponse, error) {

	// Validate that the Unit is oneof the enum
	err := s.validation.Validate.Var(in.GetUnit(), "oneof=i is cup tbsp tsp g kg ml l")
	if err != nil {
		logger.WithError(err).Errorf("Error validating unit %v", in.GetUnit())
		return nil, err
	}
	res, err := api.ConvertToBaseUnitFromRequest(in.Amount, (messages.UnitRequest)(in.Unit))

	if err != nil {
		logger.WithError(err).Error("Error converting to base unit")
		return nil, err
	}

	ingredient, err := s.dbh.InsertOneUserInventory(logger, &db.UserInventory{
		UserID:       in.GetUserId(),
		IngredientID: in.GetId(),
		Quantity:     in.GetAmount(),
		Unit:         res.Unit,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	if err != nil {
		logger.WithError(err).Error("Error inserting ingredient")
		return nil, err
	}

	return NewIngredientResponse(ingredient), nil
}

func (s *server) UpdateIngredient(ctx context.Context, in *pb.PostIngredientRequest) (*pb.IngredientResponse, error) {
	err := s.validation.Validate.Var(in.GetUnit(), "oneof=i is cup tbsp tsp g kg ml l")
	if err != nil {
		logger.WithError(err).Errorf("Error validating unit %v", in.GetUnit())
		return nil, err
	}

	res, err := api.ConvertToBaseUnitFromRequest(in.Amount, (messages.UnitRequest)(in.Unit))

	if err != nil {
		logger.WithError(err).Error("Error converting to base unit")
		return nil, err
	}

	ingredient, err := s.dbh.UpdateOneUserInventory(logger, &db.UserInventory{
		UserID:       in.GetUserId(),
		IngredientID: in.GetId(),
		Quantity:     in.GetAmount(),
		Unit:         res.Unit,
		UpdatedAt:    time.Now(),
	})

	if err != nil {
		logger.WithError(err).Error("Error updating inventory")
		return nil, err
	}

	return NewIngredientResponse(ingredient), nil
}

func (s *server) DeleteIngredient(ctx context.Context, in *pb.DeleteIngredientRequest) (*emptypb.Empty, error) {
	err := s.dbh.DeleteOneUserInventory(logger, in.GetUserId(), in.GetIngredientId())
	if err != nil {
		logger.WithError(err).Error("Error deleting inventory")
		return nil, err
	}
	return nil, nil
}

func NewServer(dbh db.DbHandler, conf *configuration.Configuration) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterInventoryServer(s, &server{
		dbh: dbh,
		// TODO MOdify the validation process with buf and validate-go
		validation: validation.New(conf),
	})
	return s
}

func Run(dbh db.DbHandler, port int, conf *configuration.Configuration) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	s := NewServer(dbh, conf)
	return s.Serve(lis)
}
