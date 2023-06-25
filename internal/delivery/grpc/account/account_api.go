package account

import (
	"context"

	"github.com/f0rmul/payment-service/internal/domain/entity"
	"github.com/f0rmul/payment-service/pkg/logger"
	pb "github.com/f0rmul/payment-service/pkg/payment-service"
)

type AccountUsecase interface {
	CreateAccount(context.Context, *entity.Account) (string, error)
}

type AccountServiceAPI struct {
	accountrUsecase AccountUsecase
	logger          logger.Logger
	pb.UnimplementedAccountServiceServer
}

func NewAccountServiceAPI(usecase AccountUsecase, logger logger.Logger) *AccountServiceAPI {
	return &AccountServiceAPI{accountrUsecase: usecase, logger: logger}
}
