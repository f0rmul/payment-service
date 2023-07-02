package account

import (
	"context"

	pb "github.com/f0rmul/payment-service/pkg/payment-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccountServiceAPI) GetAccountBalance(ctx context.Context, req *pb.AccountBalanceRequest) (*pb.AccountBalanceResponse, error) {
	balance, err := s.accountUsecase.GetBalanceByID(ctx, req.GetId())

	if err != nil {
		s.logger.Errorf("accountrUsecase.GetBalanceByID(): %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.AccountBalanceResponse{
		Id:      req.GetId(),
		Balance: uint64(balance.Balance),
	}, nil
}
