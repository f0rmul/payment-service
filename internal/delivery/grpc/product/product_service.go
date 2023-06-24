package product

import (
	"github.com/f0rmul/payment-service/pkg/logger"
	pb "github.com/f0rmul/payment-service/pkg/payment-service"
)

type ProductService struct {
	logger logger.Logger
	pb.UnimplementedProductServiceServer
}

func NewProductService(logger logger.Logger) *ProductService {
	return &ProductService{logger: logger}
}
