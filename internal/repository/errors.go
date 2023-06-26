package repository

import "errors"

var ErrProductNotFound = errors.New("product not found")
var ErrAccountNotFound = errors.New("account no found")

var ErrProductExists = errors.New("product exists")
var ErrAccountExists = errors.New("account exists")
var ErrNotEnoughMany = errors.New("not enough many on balance")
