package entity

import "auth/internal/domain"

type Inventory struct {
	Type     string `db:"type"`
	Quantity int    `db:"quantity"`
}

func (Inventory) FromDomain(i domain.Inventory) Inventory {
	return Inventory{
		Type:     i.Type,
		Quantity: i.Quantity,
	}
}

func (i Inventory) ToDomain() domain.Inventory {
	return domain.Inventory{
		Type:     i.Type,
		Quantity: i.Quantity,
	}
}

func (Inventory) SliceToDomain(slice []Inventory) []domain.Inventory {
	res := make([]domain.Inventory, len(slice))
	for i, v := range slice {
		res[i] = v.ToDomain()
	}
	return res
}
