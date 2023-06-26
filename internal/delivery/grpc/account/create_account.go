package account

import (
	"context"
	"errors"

	"github.com/f0rmul/payment-service/internal/domain"
	"github.com/f0rmul/payment-service/internal/domain/entity"
	pb "github.com/f0rmul/payment-service/pkg/payment-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccountServiceAPI) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	accountID, err := s.accountUsecase.CreateAccount(ctx, &entity.Account{
		OwnerID: req.GetOwnerId(),
	})

	if err != nil {
		s.logger.Errorf("accountrUsecase.CreateAccount(): %v", err)
		if errors.Is(err, domain.ErrAccountAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateAccountResponse{
		Id: accountID,
	}, nil
}
