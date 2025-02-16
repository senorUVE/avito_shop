package repository

import (
	"context"

	"github.com/WantBeASleep/goooool/daolib"
	"github.com/jmoiron/sqlx"
)

type DAO interface {
	daolib.DAO
	NewUserQuery(ctx context.Context) UserQuery
	NewInfoQuery(ctx context.Context) InfoQuery
	NewTransQuery(ctx context.Context) TransQuery
	NewInvQuery(ctx context.Context) InvQuery
}

type dao struct {
	daolib.DAO
}

func NewRepository(psql *sqlx.DB) DAO {
	return &dao{DAO: daolib.NewDao(psql)}
}

func (d *dao) NewUserQuery(ctx context.Context) UserQuery {
	userQuery := &userQuery{}
	d.NewRepo(ctx, userQuery)

	return userQuery
}

// /////
func (d *dao) NewInfoQuery(ctx context.Context) InfoQuery {
	infoQuery := &infoQuery{}
	d.NewRepo(ctx, infoQuery)

	return infoQuery
}

func (d *dao) NewTransQuery(ctx context.Context) TransQuery {
	transQuery := &transQuery{}
	d.NewRepo(ctx, transQuery)

	return transQuery
}

func (d *dao) NewInvQuery(ctx context.Context) InvQuery {
	invQuery := &invQuery{}
	d.NewRepo(ctx, invQuery)

	return invQuery
}
