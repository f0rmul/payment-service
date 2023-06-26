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

const accountTable = "accounts"

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, acc *entity.Account) (*entity.Account, error) {
	query, args, _ := postgres.StatementBuilder.
		Insert(accountTable).
		Columns("owner_id").
		Values(acc.OwnerID).Suffix("RETURNING *").ToSql()

	createdAccount := new(entity.Account)

	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(createdAccount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrAccountExists
		}
		return nil, errors.Wrap(err, "account.db.QueryRowxContex()")
	}
	return createdAccount, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, accountID string) (*entity.Account, error) {

	query, args, _ := postgres.StatementBuilder.
		Select("*").
		From(accountTable).Where(sq.Eq{"id": accountID}).ToSql()

	account := new(entity.Account)

	err := r.db.Get(account, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrAccountNotFound
		}
		return nil, errors.Wrap(err, "account.db.Get()")
	}

	return account, nil
}

func (r *AccountRepository) GetBalanceByID(ctx context.Context, accountID string) (int64, error) {
	query, args, _ := postgres.StatementBuilder.
		Select("balance").From(accountTable).
		Where(sq.Eq{"id": accountID}).ToSql()

	var balance int64
	err := r.db.Get(&balance, query, args...)

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

	_, err := r.db.ExecContext(ctx, query, args...)

	if err != nil {
		return errors.Wrap(err, "account.db.ExecContext()")
	}
	return nil
}

func (r *AccountRepository) Withdraw(ctx context.Context, accountID string, amount int64) error {

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
		return errors.Wrap(err, "account.db.BeginTxx()")
	}

	defer func() { _ = tx.Rollback() }()

	enoughManyQuery, args, _ := postgres.StatementBuilder.
		Select("balance").From(accountTable).
		Where(sq.Eq{"id": accountID}).ToSql()

	var balance int64

	err = tx.Get(&balance, enoughManyQuery, args...)

	if err != nil {
		return errors.Wrap(err, "account.tx.Get()")
	}

	if balance < amount {
		return repository.ErrNotEnoughMany
	}

	updateBalanceQuery, args, _ := postgres.StatementBuilder.
		Update(accountTable).Set("balance", sq.Expr("balance - ?", amount)).
		Where(sq.Eq{"id": accountID}).ToSql()

	_, err = tx.ExecContext(ctx, updateBalanceQuery, args...)

	if err != nil {
		return errors.Wrap(err, "account.tx.ExecContext()")
	}

	err = tx.Commit()

	if err != nil {
		return errors.Wrap(err, "account.tx.Commit()")
	}
	return nil
}

func (r *AccountRepository) Transfer(ctx context.Context, accountFrom, accountTo string, amount int64) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
		return errors.Wrap(err, "account.db.BeginTxx()")
	}

	defer func() { _ = tx.Rollback() }()

	enoughManyFromQuery, args, _ := postgres.StatementBuilder.
		Select("balance").From(accountTable).
		Where(sq.Eq{"id": accountFrom}).ToSql()

	var balance int64

	err = tx.Get(&balance, enoughManyFromQuery, args...)

	if err != nil {
		return errors.Wrap(err, "account.tx.Get()")
	}

	if balance < amount {
		return repository.ErrNotEnoughMany
	}

	updateBalanceFromQuery, args, _ := postgres.StatementBuilder.
		Update(accountTable).Set("balance", sq.Expr("balance - ?", amount)).
		Where(sq.Eq{"id": accountFrom}).ToSql()

	_, err = tx.ExecContext(ctx, updateBalanceFromQuery, args...)

	if err != nil {
		return errors.Wrap(err, "account.tx.ExecContext()")
	}

	updateBalanceToQuery, args, _ := postgres.StatementBuilder.
		Update(accountTable).Set("balance", sq.Expr("balance + ?", amount)).
		Where(sq.Eq{"id": accountTo}).ToSql()

	_, err = tx.ExecContext(ctx, updateBalanceToQuery, args...)

	if err != nil {
		return errors.Wrap(err, "account.tx.ExecContext()")
	}

	err = tx.Commit()

	if err != nil {
		return errors.Wrap(err, "account.tx.Commit()")
	}
	return nil
}
