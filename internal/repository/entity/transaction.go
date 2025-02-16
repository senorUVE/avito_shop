package entity

import (
	"auth/internal/domain"

	"github.com/google/uuid"
)

type Transaction struct {
	FromUser uuid.UUID `db:"from_user"`
	ToUser   uuid.UUID `db:"to_user"`
	Amount   int       `db:"amount"`
}

func (Transaction) FromDomain(t domain.Transaction) Transaction {
	return Transaction{
		FromUser: t.FromUser,
		ToUser:   t.ToUser,
		Amount:   t.Amount,
	}
}

func (t Transaction) ToDomain() domain.Transaction {
	return domain.Transaction{
		FromUser: t.FromUser,
		ToUser:   t.ToUser,
		Amount:   t.Amount,
	}
}

func (Transaction) SliceToDomain(slice []Transaction) []domain.Transaction {
	res := make([]domain.Transaction, len(slice))
	for i, v := range slice {
		res[i] = v.ToDomain()
	}
	return res
}
