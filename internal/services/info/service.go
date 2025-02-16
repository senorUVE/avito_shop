package info

import (
	"context"
	"errors"
	"log/slog"

	"auth/internal/domain"
	"auth/internal/repository"
	"auth/internal/repository/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	GetInfo(ctx context.Context) (domain.Info, error)
}

type service struct {
	dao repository.DAO
}

func New(
	dao repository.DAO,
) Service {
	return &service{
		dao: dao,
	}
}

func (s *service) GetInfo(ctx context.Context) (domain.Info, error) {

	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		slog.Error("Context is not gin.Context")
		return domain.Info{}, errors.New("context is not gin.Context")
	}

	userIdRaw, exists := ginCtx.Get("userId")
	if !exists {
		slog.Error("Unauthorized: userId not found in context")
		return domain.Info{}, errors.New("unauthorized: userId not found in context")
	}

	userId, ok := userIdRaw.(uuid.UUID)
	if !ok {
		slog.Error("Invalid userId format", "userIdRaw", userIdRaw)
		return domain.Info{}, errors.New("invalid userId format")
	}

	infoQuery := s.dao.NewInfoQuery(ctx)
	userInfo, err := infoQuery.GetUserInfo(userId)
	if err != nil {
		slog.Error("Failed to retrieve user info", "userId", userId, "error", err)
		return domain.Info{}, err
	}

	transQuery := s.dao.NewTransQuery(ctx)
	sentTransactions, err := transQuery.GetSentTransactions(userId)
	if err != nil {
		slog.Error("Failed to retrieve sent transactions", "userId", userId, "error", err)
		return domain.Info{}, err
	}
	receivedTransactions, err := transQuery.GetReceivedTransactions(userId)
	if err != nil {
		slog.Error("Failed to retrieve received transactions", "userId", userId, "error", err)
		return domain.Info{}, err
	}

	invQuery := s.dao.NewInvQuery(ctx)
	inventory, err := invQuery.GetUserInventory(userId)
	if err != nil {
		slog.Error("Failed to retrieve user inventory", "userId", userId, "error", err)
		return domain.Info{}, err
	}
	info := domain.Info{
		UserId:    userId,
		Coins:     userInfo.Coins,
		Inventory: entity.Inventory{}.SliceToDomain(inventory),
		CoinHistory: domain.CoinHistory{
			Sent:     entity.Transaction{}.SliceToDomain(sentTransactions),
			Received: entity.Transaction{}.SliceToDomain(receivedTransactions),
		},
	}
	slog.Info("Successfully retrieved user info", "userId", userId, "coins", userInfo.Coins)
	return info, nil
}
