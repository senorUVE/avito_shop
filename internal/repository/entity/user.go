package entity

import (
	"database/sql"

	"auth/internal/domain"

	"github.com/WantBeASleep/goooool/gtclib"
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
	// Coins    int            `db:"coins"`
	Token sql.NullString `db:"token"`
}

func (e User) ToDomain() domain.User {
	return domain.User{
		Id:       e.Id,
		Username: e.Username,
		Password: e.Password,
		// Coins:    e.Coins,
		Token: gtclib.String.SqlToPointer(e.Token),
	}
}

func (User) FromDomain(d domain.User) User {
	return User{
		Id:       d.Id,
		Username: d.Username,
		Password: d.Password,
		// Coins:    d.Coins,
		Token: gtclib.String.PointerToSql(d.Token),
	}
}
