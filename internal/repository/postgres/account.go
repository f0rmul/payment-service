package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/f0rmul/payment-service/internal/domain/entity"
	"github.com/f0rmul/payment-service/internal/repository"
	"github.com/f0rmul/payment-service/pkg/postgres"
)

const accountTable = "accounts"

type AccountRepository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewAccountRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) *AccountRepository {
	return &AccountRepository{db: db, getter: getter}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, acc *entity.Account) (*entity.Account, error) {
	query, args, _ := postgres.StatementBuilder.
		Insert(accountTable).
		Columns("owner_id").
		Values(acc.OwnerID).Suffix("RETURNING *").ToSql()

	createdAccount := new(repository.Account)

	err := r.getter.DefaultTrOrDB(ctx, r.db).
		QueryRowxContext(ctx, query, args...).StructScan(createdAccount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrAlreadyExists
		}
		return nil, errors.Wrap(err, "account.db.QueryRowxContex()")
	}
	return createdAccount.ToEntity(), nil
}

func (r *AccountRepository) GetByID(ctx context.Context, accountID string) (*entity.Account, error) {

	query, args, _ := postgres.StatementBuilder.
		Select("*").
		From(accountTable).Where(sq.Eq{"id": accountID}).ToSql()

	account := new(repository.Account)

	err := r.getter.DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, account, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, errors.Wrap(err, "account.db.GetContext()")
	}

	return account.ToEntity(), nil
}

func (r *AccountRepository) GetBalanceByID(ctx context.Context, accountID string) (int64, error) {
	query, args, _ := postgres.StatementBuilder.
		Select("balance").From(accountTable).
		Where(sq.Eq{"id": accountID}).ToSql()

	var balance int64
	err := r.getter.DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, &balance, query, args...)

	if err != nil {
		return 0, errors.Wrap(err, "account.db.Get()")
	}

	return balance, nil
}

func (r *AccountRepository) Deposit(ctx context.Context, accountID string, amount int64) error {

	query, args, _ := postgres.StatementBuilder.
		Update(accountTable).
		Set("balance", sq.Expr("balance + ?", amount)).
		Where(sq.Eq{"id": accountID}).ToSql()

	_, err := r.getter.DefaultTrOrDB(ctx, r.db).
		ExecContext(ctx, query, args...)

	if err != nil {
		return errors.Wrap(err, "account.db.ExecContext()")
	}
	return nil
}

func (r *AccountRepository) Withdraw(ctx context.Context, accountID string, amount int64) error {

	tx := r.getter.DefaultTrOrDB(ctx, r.db)

	updateBalanceQuery, args, _ := postgres.StatementBuilder.
		Update(accountTable).Set("balance", sq.Expr("balance - ?", amount)).
		Where(sq.Eq{"id": accountID}).ToSql()

	_, err := tx.ExecContext(ctx, updateBalanceQuery, args...)

	if err != nil {
		return errors.Wrap(err, "account.tx.ExecContext()")
	}
	return nil
}
