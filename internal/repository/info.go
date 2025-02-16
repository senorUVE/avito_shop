package repository

import (
	"auth/internal/repository/entity"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/WantBeASleep/goooool/daolib"
	"github.com/google/uuid"
)

const infoTable = "\"info\""

type InfoQuery interface {
	UpdateCoins(id uuid.UUID, balance int) error
	GetUserInfo(id uuid.UUID) (entity.Info, error)
	InsertUserInfo(userId uuid.UUID, coins int) error
}

type infoQuery struct {
	*daolib.BaseQuery
}

func (q *infoQuery) SetBaseQuery(baseQuery *daolib.BaseQuery) {
	q.BaseQuery = baseQuery
}

func (q *infoQuery) UpdateCoins(id uuid.UUID, balance int) error {
	query := q.QueryBuilder().
		Update(infoTable).
		Set("coins", balance).
		Where(sq.Eq{
			"user_id": id,
		})

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}
	return nil
}

func (q *infoQuery) GetUserInfo(id uuid.UUID) (entity.Info, error) {
	query := q.QueryBuilder().
		Select(
			"user_id",
			"coins",
		).
		From(infoTable).
		Where(sq.Eq{
			"user_id": id,
		})

	sqlStr, args, _ := query.ToSql()
	slog.Info("Executing query for GetUserInfo", "query", sqlStr, "args", args)
	var userInfo entity.Info
	err := q.Runner().Getx(q.Context(), &userInfo, query)
	if err != nil {
		slog.Error("Failed to execute GetUserInfo query", "userId", id, "error", err)
		return entity.Info{}, err
	}
	return userInfo, err
}

func (q *infoQuery) InsertUserInfo(userId uuid.UUID, coins int) error {
	query := q.QueryBuilder().
		Insert("info").
		Columns(
			"user_id",
			"coins",
		).
		Values(userId, coins)

	_, err := q.Runner().Execx(q.Context(), query)
	return err
}
