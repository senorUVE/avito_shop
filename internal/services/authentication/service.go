package authentication

import (
	"context"
	"errors"
	"fmt"

	"auth/internal/repository"
	"auth/internal/repository/entity"
	"auth/internal/services/password"
	"auth/internal/services/tokenizer"

	"github.com/google/uuid"
)

type Service interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type service struct {
	dao          repository.DAO
	tokenizerSrv tokenizer.Service
	passwordSrv  password.Service
}

func New(
	dao repository.DAO,
	tokenizerSrv tokenizer.Service,
	passwordSrv password.Service,
) Service {
	return &service{
		dao:          dao,
		tokenizerSrv: tokenizerSrv,
		passwordSrv:  passwordSrv,
	}
}

func (s *service) Authenticate(ctx context.Context, username, password string) (string, error) {
	userQuery := s.dao.NewUserQuery(ctx)
	userDB, err := userQuery.GetUserByUsername(username)
	if err != nil {
		salt := s.passwordSrv.GenerateSalt(ctx)
		hashedPassword, err := s.passwordSrv.Hash(ctx, password, salt)
		if err != nil {
			return "", fmt.Errorf("hash password: %w", err)
		}

		newUser := entity.User{
			Id:       uuid.New(),
			Username: username,
			Password: hashedPassword,
		}
		if err := userQuery.InsertUser(newUser); err != nil {
			return "", fmt.Errorf("create user: %w", err)
		}
		infoQuery := s.dao.NewInfoQuery(ctx)
		if err := infoQuery.InsertUserInfo(newUser.Id, 1000); err != nil {
			return "", fmt.Errorf("initialize user balance: %w", err)
		}

		userDB = newUser
	} else {
		salt, err := s.passwordSrv.GetSalt(ctx, userDB.Password)
		if err != nil {
			return "", fmt.Errorf("get salt: %w", err)
		}
		hashedPassword, err := s.passwordSrv.Hash(ctx, password, salt)
		if err != nil {
			return "", fmt.Errorf("hash password: %w", err)
		}
		if hashedPassword != userDB.Password {
			return "", errors.New("invalid password")
		}
	}
	user := userDB.ToDomain()

	accessToken, refreshToken, err := s.tokenizerSrv.GeneratePair(
		ctx,
		map[string]any{"x-user_id": user.Id.String()},
	)
	if err != nil {
		return "", fmt.Errorf("generate tokens: %w", err)
	}
	user.Token = &refreshToken
	if _, err := userQuery.UpdateUser(entity.User{}.FromDomain(user)); err != nil {
		return "", fmt.Errorf("update user token: %w", err)
	}

	return accessToken, nil
}
