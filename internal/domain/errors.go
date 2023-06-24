package domain

import "errors"

var ErrProductAlreadyExists = errors.New("product already exists")
var ErrProductNotFound = errors.New("product not found")
