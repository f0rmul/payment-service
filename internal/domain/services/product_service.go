package services

import (
	"context"

	"github.com/f0rmul/payment-service/internal/domain/entity"
	"github.com/f0rmul/payment-service/pkg/logger"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product entity.Product) (*entity.Product, error)
	GetByID(ctx context.Context, productID string) (*entity.Product, error)
}

type ProductService struct {
	repo   ProductRepository
	logger logger.Logger
}

func NewProductService(r ProductRepository, logger logger.Logger) *ProductService {
	return &ProductService{repo: r, logger: logger}
}
