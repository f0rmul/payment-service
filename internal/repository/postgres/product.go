package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/f0rmul/payment-service/internal/domain/entity"
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

	query, args, err := postgres.StatementBuilder.Insert(productTable).
		Columns("id", "name").Values(product.ID, product.Name).
		Suffix("RETURNING *").ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "query builder error")
	}

	createdProduct := new(entity.Product)
	err = r.db.QueryRowxContext(ctx, query, args...).StructScan(createdProduct)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "db.QueryRowxContex()")
	}
	return createdProduct, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, productID string) (*entity.Product, error) {

	query, args, err := postgres.StatementBuilder.Select("*").
		From(productTable).
		Where(sq.Eq{"id": productID}).ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "query builder error")
	}

	product := new(entity.Product)
	err = r.db.QueryRowxContext(ctx, query, args...).StructScan(product)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "db.QueryRowxContex()")
	}

	return product, nil
}
