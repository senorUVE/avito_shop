package repository

import (
	"auth/internal/repository/entity"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/WantBeASleep/goooool/daolib"
	"github.com/google/uuid"
)

const transTable = "\ncoin_transactions\n"

type TransQuery interface {
	InsertTransaction(fromUser, toUser uuid.UUID, amount int) error
	GetSentTransactions(userId uuid.UUID) ([]entity.Transaction, error)
	GetReceivedTransactions(userId uuid.UUID) ([]entity.Transaction, error)
}

type transQuery struct {
	*daolib.BaseQuery
}

func (q *transQuery) SetBaseQuery(baseQuery *daolib.BaseQuery) {
	q.BaseQuery = baseQuery
}

func (q *transQuery) InsertTransaction(fromUser, toUser uuid.UUID, amount int) error {
	query := q.QueryBuilder().
		Insert(
			"coin_transactions",
		).
		Columns(
			"from_user",
			"to_user",
			"amount",
		).
		Values(fromUser, toUser, amount)

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}
	return nil
}

func (q *transQuery) GetSentTransactions(userId uuid.UUID) ([]entity.Transaction, error) {
	query := q.QueryBuilder().
		Select(
			"from_user",
			"to_user",
			"amount",
		).
		From(transTable).
		Where(sq.Eq{
			"from_user": userId,
		})

	fmt.Println("Executing query on table:", transTable)
	var transactions []entity.Transaction
	err := q.Runner().Selectx(q.Context(), &transactions, query)
	return transactions, err
}

func (q *transQuery) GetReceivedTransactions(userId uuid.UUID) ([]entity.Transaction, error) {
	query := q.QueryBuilder().
		Select(
			"from_user",
			"to_user",
			"amount",
		).
		From(transTable).
		Where(sq.Eq{
			"to_user": userId,
		})

	var transactions []entity.Transaction
	err := q.Runner().Selectx(q.Context(), &transactions, query)
	return transactions, err
}
