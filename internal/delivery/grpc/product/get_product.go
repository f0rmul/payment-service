package product

import (
	"context"

	pb "github.com/f0rmul/payment-service/pkg/payment-service"
)

func (s *ProductService) GetProductByID(ctx context.Context, request *pb.ProductByIDRequest) (*pb.ProductByIDResponse, error) {
	panic("unimplemented")
}
