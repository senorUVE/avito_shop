package domain

import (
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	Username string
	Password string
	// Coins    int
	Token *string
}

// type Merch struct {
// 	Name  string
// 	Price int
// }

// type Inventory struct {
// 	Type     string
// 	Quantity int
// }
