package product

import (
	"context"
	"errors"

	"github.com/f0rmul/payment-service/internal/domain"
	pb "github.com/f0rmul/payment-service/pkg/payment-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ProductServiceAPI) GetProductByID(ctx context.Context, request *pb.ProductByIDRequest) (*pb.ProductByIDResponse, error) {

	product, err := s.productUsecase.GetProductByID(ctx, request.GetId())

	if err != nil {
		s.logger.Errorf("productUsecase.GetProductByID(): %v", err)
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ProductByIDResponse{
		Product: &pb.Product{
			Id:   product.ID,
			Name: product.Name,
		},
	}, nil
}
