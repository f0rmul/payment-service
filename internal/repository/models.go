package repository

import (
	"time"

	"github.com/f0rmul/payment-service/internal/domain/entity"
)

type Account struct {
	ID        string    `db:"id"`
	OwnerID   string    `db:"owner_id"`
	Balance   int64     `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
}

func (a *Account) ToEntity() *entity.Account {
	return &entity.Account{
		ID:        a.ID,
		OwnerID:   a.OwnerID,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
	}
}
