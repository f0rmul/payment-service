package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type WithTxReturnError func(context.Context, *sqlx.Tx) error

func WithTransactionReturnError(ctx context.Context, tFunc WithTxReturnError) error {
	panic("unimplemented")
}
