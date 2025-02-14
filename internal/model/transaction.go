package model

type Transaction struct {
	ID         string
	FromUserID string
	ToUserID   string
	Amount     int
}
