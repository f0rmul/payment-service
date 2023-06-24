package entity

import "time"

type Account struct {
	ID        string    `db:"id"`
	OwnerID   string    `db:"owner_id"`
	Balance   int64     `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
}
