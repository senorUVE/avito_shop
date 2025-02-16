package domain

import "github.com/google/uuid"

type Inventory struct {
	Type     string
	Quantity int
}

type Info struct {
	UserId      uuid.UUID
	Coins       int
	Inventory   []Inventory
	CoinHistory CoinHistory
}

type Transaction struct {
	FromUser uuid.UUID
	ToUser   uuid.UUID
	Amount   int
}

type CoinHistory struct {
	Received []Transaction
	Sent     []Transaction
}
