package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/f0rmul/payment-service/internal/domain/entity"
	"github.com/f0rmul/payment-service/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func PrepareMockDB(t *testing.T) (*sqlx.DB, *sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return sqlxDB, db, mock
}
func TestCreateAccount(t *testing.T) {
	t.Parallel()

	sqlxDB, db, mock := PrepareMockDB(t)
	defer sqlxDB.Close()
	defer db.Close()

	accountRepo := NewAccountRepository(sqlxDB, trmsqlx.DefaultCtxGetter)

	columns := []string{"id", "owner_id", "balance", "created_at"}

	accountID := uuid.New()
	ownerID := uuid.New()
	mockAccount := entity.Account{
		OwnerID: ownerID.String(),
	}

	rows := sqlmock.NewRows(columns).AddRow(
		accountID,
		mockAccount.OwnerID,
		0,
		time.Now())

	query, _, _ := postgres.StatementBuilder.
		Insert(accountTable).
		Columns("owner_id").
		Values(mockAccount.OwnerID).Suffix("RETURNING *").ToSql()

	mock.ExpectQuery(query).WithArgs(
		mockAccount.OwnerID,
	).WillReturnRows(rows)

	createdAcc, err := accountRepo.CreateAccount(context.TODO(), &mockAccount)
	require.NoError(t, err)
	require.NotNil(t, createdAcc)
	require.Equal(t, accountID.String(), createdAcc.ID)
}

func TestAccountByID(t *testing.T) {
	t.Parallel()

	sqlxDB, db, mock := PrepareMockDB(t)
	defer sqlxDB.Close()
	defer db.Close()

	accountRepo := NewAccountRepository(sqlxDB, trmsqlx.DefaultCtxGetter)

	columns := []string{"id", "owner_id", "balance", "created_at"}

	accountID := uuid.New()
	ownerID := uuid.New()
	creationTime := time.Now()

	mockAccount := entity.Account{
		ID:        accountID.String(),
		OwnerID:   ownerID.String(),
		Balance:   1500,
		CreatedAt: creationTime,
	}

	rows := sqlmock.NewRows(columns).AddRow(
		mockAccount.ID,
		mockAccount.OwnerID,
		mockAccount.Balance,
		mockAccount.CreatedAt,
	)

	query, _, _ := postgres.StatementBuilder.
		Select("*").
		From(accountTable).Where(sq.Eq{"id": accountID}).ToSql()

	mock.ExpectQuery(query).WithArgs(
		mockAccount.ID,
	).WillReturnRows(rows)

	account, err := accountRepo.GetByID(context.TODO(), accountID.String())
	require.NoError(t, err)
	require.NotNil(t, account)

	require.Equal(t, account.Balance, mockAccount.Balance)
	require.Equal(t, account.CreatedAt, creationTime)
	require.Equal(t, account.OwnerID, ownerID.String())
}

func TestAccountDeposit(t *testing.T) {
	t.Parallel()

	sqlxDB, db, mock := PrepareMockDB(t)
	defer sqlxDB.Close()
	defer db.Close()

	accountRepo := NewAccountRepository(sqlxDB, trmsqlx.DefaultCtxGetter)
	accountID := uuid.New()
	amount := 500

	query, _, _ := postgres.StatementBuilder.
		Update(accountTable).
		Set("balance", sq.Expr("balance + ?", amount)).
		Where(sq.Eq{"id": accountID}).ToSql()

	mock.ExpectExec(query).WithArgs(
		amount,
		accountID.String(),
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err := accountRepo.Deposit(context.TODO(), accountID.String(), int64(amount))

	require.NoError(t, err)
}
