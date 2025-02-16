package model

import "errors"

var (
	ErrItemNotFound       = errors.New("item not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrInventory          = errors.New("failed to get inventory")
	ErrNegAmount          = errors.New("amount must be positive")
	ErrHistory            = errors.New("failed to get history")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalError      = errors.New("something went wrong")
	ErrCreateUser         = errors.New("create user error")
)
