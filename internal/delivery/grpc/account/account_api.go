package account

import (
	"context"

	"github.com/f0rmul/payment-service/internal/domain/entity"
	"github.com/f0rmul/payment-service/pkg/logger"
	pb "github.com/f0rmul/payment-service/pkg/payment-service"
)

type AccountUsecase interface {
	CreateAccount(context.Context, *entity.Account) (string, error)
	GetBalanceByID(context.Context, string) (*entity.Account, error)
	Credit(context.Context, string, int64) (*entity.Account, error)
	WriteOff(context.Context, string, int64) error
	Transfer(context.Context, string, string, int64) error
}

type AccountServiceAPI struct {
	accountUsecase AccountUsecase
	logger         logger.Logger
	pb.UnimplementedAccountServiceServer
}

func NewAccountServiceAPI(usecase AccountUsecase, logger logger.Logger) *AccountServiceAPI {
	return &AccountServiceAPI{accountUsecase: usecase, logger: logger}
}
