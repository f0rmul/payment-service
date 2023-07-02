package entity

import (
	"time"
)

type Account struct {
	ID        string
	OwnerID   string
	Balance   int64
	CreatedAt time.Time
}
