package product_service

import (
	"context"
	"testLmd/libs/errors"
	product_service "testLmd/services/product-service/proto/product-service"
	"testLmd/services/product-service/repository"
)

type Server struct {
	product_service.UnimplementedProductServiceServer
}

var db repository.Database

func Init() error {
	var err error

	db, err = repository.NewDatabase()
	return err
}

func (s Server) GetAllProducts(
	ctx context.Context,
	r *product_service.GetAllProductsRequest,
) (*product_service.GetAllProductsResponse, error) {
	err := r.Validate()
	if err != nil {
		return nil, errors.LogError(err)
	}

	products, err := db.GetAllProducts(ctx, r.WarehouseId)
	if err != nil {
		return nil, errors.LogError(err)
	}

	response := &product_service.GetAllProductsResponse{Product: products}

	return response, err
}

func (s Server) CancelReservationProduct(
	ctx context.Context,
	r *product_service.CancelReservationProductRequest,
) (*product_service.CancelReservationProductResponse, error) {
	err := r.Validate()
	if err != nil {
		return nil, errors.LogError(err)
	}

	out, err := db.CancelReservationProduct(ctx, r.Id, r.WarehouseId)
	if err != nil {
		return nil, errors.LogError(err)
	}

	return &product_service.CancelReservationProductResponse{ProductId: out}, err
}

func (s Server) ReservationProduct(
	ctx context.Context,
	r *product_service.ReservationProductRequest,
) (*product_service.ReservationProductResponse, error) {
	err := r.Validate()
	if err != nil {
		return nil, errors.LogError(err)
	}

	out, err := db.ReservationProduct(ctx, r.Id, r.WarehouseId)
	if err != nil {
		return nil, errors.LogError(err)
	}

	return &product_service.ReservationProductResponse{ProductId: out}, err
}
