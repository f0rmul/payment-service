package product

import (
	"context"

	"github.com/f0rmul/payment-service/internal/domain/entity"
	"github.com/f0rmul/payment-service/pkg/logger"
	pb "github.com/f0rmul/payment-service/pkg/payment-service"
)

type ProductUsecase interface {
	CreateProduct(context.Context, *entity.Product) (string, error)
	GetProductByID(context.Context, string) (*entity.Product, error)
}

type ProductServiceAPI struct {
	productUsecase ProductUsecase
	logger         logger.Logger
	pb.UnimplementedProductServiceServer
}

func NewProductService(usecase ProductUsecase, logger logger.Logger) *ProductServiceAPI {
	return &ProductServiceAPI{
		productUsecase: usecase,
		logger:         logger,
	}
}
