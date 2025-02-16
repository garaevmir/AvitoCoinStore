package model

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"required,gt=0"`
}

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
