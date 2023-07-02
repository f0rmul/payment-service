package domain

import "errors"

var ErrProductAlreadyExists = errors.New("product already exists")
var ErrProductNotFound = errors.New("product not found")

var ErrAccountAlreadyExists = errors.New("account already exists")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughMoney = errors.New("not enough money on account")
