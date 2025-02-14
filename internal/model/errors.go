package model

import "errors"

var (
	ErrItemNotFound       = errors.New("item not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidCredentials = errors.New(("invalid credentials"))
)
