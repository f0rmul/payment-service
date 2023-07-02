package services

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/pkg/errors"

	"github.com/f0rmul/payment-service/internal/domain"
	"github.com/f0rmul/payment-service/internal/domain/entity"
	"github.com/f0rmul/payment-service/internal/repository"
	"github.com/f0rmul/payment-service/pkg/logger"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, acc *entity.Account) (*entity.Account, error)
	GetByID(ctx context.Context, accountID string) (*entity.Account, error)
	GetBalanceByID(ctx context.Context, accountID string) (int64, error)
	Deposit(ctx context.Context, accountID string, amount int64) error
	Withdraw(ctx context.Context, accountID string, amount int64) error
}

type AccountService struct {
	accountRepo AccountRepository
	txManager   *manager.Manager
	logger      logger.Logger
}

func NewAccountService(repo AccountRepository, manager *manager.Manager, logger logger.Logger) *AccountService {
	return &AccountService{
		accountRepo: repo,
		txManager:   manager,
		logger:      logger,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, acc *entity.Account) (*entity.Account, error) {
	err := s.txManager.Do(ctx, func(ctx context.Context) error {

		_, err := s.accountRepo.GetByID(ctx, acc.OwnerID)

		if err != nil {
			s.logger.Errorf("s.accountRepo.GetByID(): %v", err)
			if errors.Is(err, repository.ErrAlreadyExists) {
				return domain.ErrAccountAlreadyExists
			}
			return err
		}

		acc, err = s.accountRepo.CreateAccount(ctx, acc)
		if err != nil {
			s.logger.Errorf("s.accountRepo.CreateAccount(): %v", err)
			return err
		}
		return nil

	})
	return acc, err
}

func (s *AccountService) GetByID(ctx context.Context, accountID string) (*entity.Account, error) {

	account, err := s.accountRepo.GetByID(ctx, accountID)

	if err != nil {
		s.logger.Errorf("s.accountRepo.GetByID(): %v", err)
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, errors.Wrap(err, "s.accountRepo.GetByID()")
	}
	return account, nil
}

func (s *AccountService) GetBalanceByID(ctx context.Context, accountID string) (int64, error) {

	balance, err := s.accountRepo.GetBalanceByID(ctx, accountID)

	if err != nil {
		s.logger.Errorf("s.accountRepo.GetBalanceByID(): %v", err)
		return 0, err
	}
	return balance, nil
}

func (s *AccountService) WriteOff(ctx context.Context, accountID string, amount int64) error {
	err := s.txManager.Do(ctx, func(ctx context.Context) error {

		account, err := s.accountRepo.GetByID(ctx, accountID)

		if err != nil {
			s.logger.Errorf("s.accountRepo.GetByID(): %v", err)
			if errors.Is(err, repository.ErrNotFound) {
				return domain.ErrAccountNotFound
			}
			return errors.Wrap(err, "s.accountRepo.GetByID()")
		}

		if account.Balance < amount {
			s.logger.Errorf("not enough money: %v", domain.ErrNotEnoughMoney)
			return domain.ErrNotEnoughMoney
		}

		err = s.accountRepo.Withdraw(ctx, accountID, amount)

		if err != nil {
			s.logger.Errorf("s.accountRepo.Withdraw(): %v", err)
			return err
		}

		return nil
	})
	return err
}

func (s *AccountService) Credit(ctx context.Context, accountID string, amount int64) error {
	err := s.txManager.Do(ctx, func(ctx context.Context) error {

		_, err := s.accountRepo.GetByID(ctx, accountID)

		if err != nil {
			s.logger.Errorf("s.accountRepo.GetByID(): %v", err)
			if errors.Is(err, repository.ErrNotFound) {
				return domain.ErrAccountNotFound
			}
			return errors.Wrap(err, "s.accountRepo.GetByID()")
		}

		err = s.accountRepo.Deposit(ctx, accountID, amount)

		if err != nil {
			s.logger.Errorf("s.accountRepo.Deposit(): %v", err)
			return err
		}

		return nil
	})
	return err
}

type TransferInput struct {
	AccountTo   string
	AccountFrom string
	Amount      int64
}

func (s *AccountService) Transfer(ctx context.Context, input TransferInput) error {

	err := s.txManager.Do(ctx, func(ctx context.Context) error {

		account, err := s.accountRepo.GetByID(ctx, input.AccountFrom)

		if err != nil {
			s.logger.Errorf("s.accountRepo.GetByID(): %v", err)
			if errors.Is(err, repository.ErrNotFound) {
				return domain.ErrAccountNotFound
			}
			return errors.Wrap(err, "s.accountRepo.GetByID()")
		}

		if account.Balance < input.Amount {
			s.logger.Errorf("not enough money: %v", domain.ErrNotEnoughMoney)
			return domain.ErrNotEnoughMoney
		}

		err = s.accountRepo.Withdraw(ctx, input.AccountFrom, input.Amount)

		if err != nil {
			s.logger.Errorf("s.accountRepo.Withdraw(): %v", err)
			return err
		}

		err = s.accountRepo.Deposit(ctx, input.AccountTo, input.Amount)

		if err != nil {
			s.logger.Errorf("s.accountRepo.Deposit(): %v", err)
			return err
		}
		return nil
	})

	return err
}
