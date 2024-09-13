package storage

import "errors"

var (
	ErrAccountExists    = errors.New("account with this id already exists")
	ErrAccountNotFound  = errors.New("account not found")
	ErrCurrencyNotFound = errors.New("currency not found")
)
