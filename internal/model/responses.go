package model

import "time"

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type InventoryItem struct {
	Name     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type TransactionHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	FromUser  string    `json:"fromUser"`
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type SentTransaction struct {
	ToUser    string    `json:"toUser"`
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type InfoResponse struct {
	Coins       int                `json:"coins"`
	Inventory   []InventoryItem    `json:"inventory"`
	CoinHistory TransactionHistory `json:"coinHistory"`
}
