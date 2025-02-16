package transaction

import (
	"context"
	"errors"
	"fmt"

	"auth/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	TransferCoins(ctx context.Context, toUser string, amount int) error
	BuyItem(ctx context.Context, itemType string, quantity int) error
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

func (s *service) TransferCoins(ctx context.Context, toUser string, amount int) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return errors.New("context is not gin.Context")
	}

	userIdRaw, exists := ginCtx.Get("userId")
	if !exists {
		return errors.New("unauthorized: userId not found in context")
	}

	fromUser, ok := userIdRaw.(uuid.UUID)
	if !ok {
		return errors.New("invalid userId format")
	}

	userQuery := s.dao.NewUserQuery(ctx)
	toUserUUID, err := userQuery.GetUserByUsername(toUser)
	if err != nil {
		return fmt.Errorf("receiver not found: %w", err)
	}

	if fromUser == toUserUUID.Id {
		return errors.New("cannot transfer to yourself")
	}

	infoQuery := s.dao.NewInfoQuery(ctx)
	senderInfo, err := infoQuery.GetUserInfo(fromUser)
	if err != nil {
		return fmt.Errorf("get sender info: %w", err)
	}

	if senderInfo.Coins < amount {
		return errors.New("insufficient funds")
	}

	receiverInfo, err := infoQuery.GetUserInfo(toUserUUID.Id)
	if err != nil {
		return fmt.Errorf("get receiver info: %w", err)
	}

	if err := infoQuery.UpdateCoins(fromUser, senderInfo.Coins-amount); err != nil {
		return fmt.Errorf("update sender balance: %w", err)
	}
	if err := infoQuery.UpdateCoins(toUserUUID.Id, receiverInfo.Coins+amount); err != nil {
		return fmt.Errorf("update receiver balance: %w", err)
	}

	transQuery := s.dao.NewTransQuery(ctx)
	if err := transQuery.InsertTransaction(fromUser, toUserUUID.Id, amount); err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	return nil
}

func (s *service) BuyItem(ctx context.Context, itemType string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return errors.New("context is not gin.Context")
	}

	userIdRaw, exists := ginCtx.Get("userId")
	if !exists {
		return errors.New("unauthorized: userId not found in context")
	}

	userId, ok := userIdRaw.(uuid.UUID)
	if !ok {
		return errors.New("invalid userId format")
	}

	itemPrices := map[string]int{
		"t-shirt":    80,
		"cup":        20,
		"book":       50,
		"pen":        10,
		"powerbank":  200,
		"hoody":      300,
		"umbrella":   200,
		"socks":      10,
		"wallet":     50,
		"pink-hoody": 500,
	}

	price, exists := itemPrices[itemType]
	if !exists {
		return errors.New("invalid item type")
	}

	totalCost := price * quantity

	infoQuery := s.dao.NewInfoQuery(ctx)
	userInfo, err := infoQuery.GetUserInfo(userId)
	if err != nil {
		return fmt.Errorf("get user info: %w", err)
	}

	if userInfo.Coins < totalCost {
		return errors.New("insufficient funds")
	}

	if err := infoQuery.UpdateCoins(userId, userInfo.Coins-totalCost); err != nil {
		return fmt.Errorf("update user balance: %w", err)
	}

	invQuery := s.dao.NewInvQuery(ctx)
	if err := invQuery.InsertInventory(userId, itemType, quantity); err != nil {
		return fmt.Errorf("insert inventory: %w", err)
	}

	return nil
}
