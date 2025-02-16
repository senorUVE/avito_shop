package entity

import (
	"auth/internal/domain"

	"github.com/google/uuid"
)

type Info struct {
	UserId uuid.UUID `db:"user_id"`
	Coins  int       `db:"coins"`
}

func (e Info) ToDomain(inventory []domain.Inventory) domain.Info {
	return domain.Info{
		UserId:    e.UserId,
		Coins:     e.Coins,
		Inventory: inventory,
	}
}

func (Info) FromDomain(d domain.Info) Info {
	return Info{
		UserId: d.UserId,
		Coins:  d.Coins,
	}
}
