package postgres

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Config interface {
	GetDSN() string
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxIdleTime() time.Duration
	GetConnMaxIdleLifeTime() time.Duration
}

var StatementBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func NewConn(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", cfg.GetDSN())

	if err != nil {
		return nil, errors.Wrap(err, "sqlx.Open()")
	}

	db.SetMaxOpenConns(cfg.GetMaxOpenConns())
	db.SetMaxIdleConns(cfg.GetMaxIdleConns())
	db.SetConnMaxIdleTime(cfg.GetConnMaxIdleTime())
	db.SetConnMaxLifetime(cfg.GetConnMaxIdleLifeTime())

	return db, nil
}
