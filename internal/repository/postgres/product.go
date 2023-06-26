package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/f0rmul/payment-service/internal/domain/entity"
	"github.com/f0rmul/payment-service/internal/repository"
	"github.com/f0rmul/payment-service/pkg/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const productTable = "products"

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product entity.Product) (*entity.Product, error) {

	query, args, _ := postgres.StatementBuilder.Insert(productTable).
		Columns("id", "name").Values(product.ID, product.Name).
		Suffix("RETURNING *").ToSql()

	createdProduct := new(entity.Product)
	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(createdProduct)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrProductExists
		}
		return nil, errors.Wrap(err, "product.db.QueryRowxContex()")
	}
	return createdProduct, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, productID string) (*entity.Product, error) {
	query, args, _ := postgres.StatementBuilder.Select("*").
		From(productTable).
		Where(sq.Eq{"id": productID}).ToSql()

	product := new(entity.Product)
	err := r.db.Get(product, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrProductNotFound
		}
		return nil, errors.Wrap(err, "product.db.Get()")
	}

	return product, nil
}
