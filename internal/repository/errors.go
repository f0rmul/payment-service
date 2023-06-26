package repository

import "errors"

var ErrNotFound = errors.New("item no found")
var ErrAlreadyExists = errors.New("item already exists")

var ErrNotEnoughMany = errors.New("not enough many on balance")
