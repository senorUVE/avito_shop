package repository

import (
	"auth/internal/repository/entity"

	sq "github.com/Masterminds/squirrel"
	"github.com/WantBeASleep/goooool/daolib"
	"github.com/google/uuid"
)

const userTable = "\"user\""

type UserQuery interface {
	InsertUser(user entity.User) error
	GetUserByPK(id uuid.UUID) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	UpdateUser(user entity.User) (int64, error)
	GetUserIdByUsername(username string) (uuid.UUID, error)
}

type userQuery struct {
	*daolib.BaseQuery
}

func (q *userQuery) SetBaseQuery(baseQuery *daolib.BaseQuery) {
	q.BaseQuery = baseQuery
}

func (q *userQuery) InsertUser(user entity.User) error {
	query := q.QueryBuilder().
		Insert(userTable).
		Columns(
			"id",
			"username",
			"password",
			//"coins",
			"token",
		).
		Values(
			user.Id,
			user.Username,
			user.Password,
			user.Token,
		)

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}

func (q *userQuery) GetUserByPK(id uuid.UUID) (entity.User, error) {
	query := q.QueryBuilder().
		Select(
			"id",
			"username",
			"password",
			//"coins",
			"token",
		).
		From(userTable).
		Where(sq.Eq{
			"id": id,
		})

	var user entity.User
	if err := q.Runner().Getx(q.Context(), &user, query); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (q *userQuery) GetUserByUsername(username string) (entity.User, error) {
	query := q.QueryBuilder().
		Select(
			"id",
			"username",
			"password",
			//"coins",
			"token",
		).
		From(userTable).
		Where(sq.Eq{
			"username": username,
		})

	var user entity.User
	if err := q.Runner().Getx(q.Context(), &user, query); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (q *userQuery) UpdateUser(user entity.User) (int64, error) {
	query := q.QueryBuilder().
		Update(userTable).
		SetMap(sq.Eq{
			"password": user.Password,
			"token":    user.Token,
		}).
		Where(sq.Eq{
			"id": user.Id,
		})

	res, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (q *userQuery) GetUserIdByUsername(username string) (uuid.UUID, error) {
	query := q.QueryBuilder().
		Select("id").
		From(userTable).
		Where(sq.Eq{
			"username": username,
		})

	var userId uuid.UUID
	err := q.Runner().Getx(q.Context(), &userId, query)
	return userId, err
}
