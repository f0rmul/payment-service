package product

import (
	"context"
	"errors"

	"github.com/f0rmul/payment-service/internal/domain"
	"github.com/f0rmul/payment-service/internal/domain/entity"
	pb "github.com/f0rmul/payment-service/pkg/payment-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ProductServiceAPI) CreateProduct(ctx context.Context, request *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {

	productID, err := s.productUsecase.CreateProduct(ctx, entity.Product{
		Name: request.GetName(),
	})

	if err != nil {
		s.logger.Errorf("productUsecase.CreateProduct(): %v", err)
		if errors.Is(err, domain.ErrProductAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateProductResponse{
		Id: productID,
	}, nil
}
