package repository

import (
	"auth/internal/repository/entity"

	sq "github.com/Masterminds/squirrel"
	"github.com/WantBeASleep/goooool/daolib"
	"github.com/google/uuid"
)

const invTable = "\"inventory\""

type InvQuery interface {
	InsertInventory(userId uuid.UUID, itemType string, quantity int) error
	GetUserInventory(userId uuid.UUID) ([]entity.Inventory, error)
}

type invQuery struct {
	*daolib.BaseQuery
}

func (q *invQuery) SetBaseQuery(baseQuery *daolib.BaseQuery) {
	q.BaseQuery = baseQuery
}

func (q *invQuery) InsertInventory(userId uuid.UUID, itemType string, quantity int) error {
	query := q.QueryBuilder().
		Insert("inventory").
		Columns(
			"user_id",
			"type",
			"quantity",
		).
		Values(
			userId,
			itemType,
			quantity,
		)

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}
	return nil
}

func (q *invQuery) GetUserInventory(userId uuid.UUID) ([]entity.Inventory, error) {
	query := q.QueryBuilder().
		Select(
			"type",
			"quantity",
		).
		From(invTable).
		Where(sq.Eq{
			"user_id": userId,
		})

	var inventory []entity.Inventory
	err := q.Runner().Selectx(q.Context(), &inventory, query)
	return inventory, err
}
