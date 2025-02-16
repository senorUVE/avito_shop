package transaction_test

import (
	"auth/internal/repository"
	"auth/internal/repository/entity"
	"auth/internal/services/transaction"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_TransferCoins_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockUserQuery := repository.NewMockUserQuery(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)
	mockTransQuery := repository.NewMockTransQuery(ctrl)

	mockDAO.EXPECT().NewUserQuery(gomock.Any()).Return(mockUserQuery)
	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery)
	mockDAO.EXPECT().NewTransQuery(gomock.Any()).Return(mockTransQuery)

	service := transaction.New(mockDAO)

	ginCtx, _ := gin.CreateTestContext(nil)
	fromUser := uuid.New()
	toUser := uuid.New()
	amount := 50
	ginCtx.Set("userId", fromUser)

	mockUserQuery.EXPECT().GetUserByUsername(gomock.Any()).Return(entity.User{
		Id:       toUser,
		Username: "receiverUser",
	}, nil)

	mockInfoQuery.EXPECT().GetUserInfo(fromUser).Return(entity.Info{Coins: 100}, nil)
	mockInfoQuery.EXPECT().GetUserInfo(toUser).Return(entity.Info{Coins: 30}, nil)

	mockInfoQuery.EXPECT().UpdateCoins(fromUser, 50).Return(nil)
	mockInfoQuery.EXPECT().UpdateCoins(toUser, 80).Return(nil)

	mockTransQuery.EXPECT().InsertTransaction(fromUser, toUser, amount).Return(nil)

	err := service.TransferCoins(ginCtx, toUser.String(), amount)
	assert.NoError(t, err)
}

func TestService_TransferCoins_InsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockUserQuery := repository.NewMockUserQuery(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)

	mockDAO.EXPECT().NewUserQuery(gomock.Any()).Return(mockUserQuery)
	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery)

	service := transaction.New(mockDAO)

	ginCtx, _ := gin.CreateTestContext(nil)
	fromUser := uuid.New()
	toUser := uuid.New()
	amount := 50
	ginCtx.Set("userId", fromUser)

	mockUserQuery.EXPECT().GetUserByUsername(toUser.String()).Return(entity.User{Id: toUser}, nil)

	mockInfoQuery.EXPECT().GetUserInfo(fromUser).Return(entity.Info{Coins: 10}, nil)

	err := service.TransferCoins(ginCtx, toUser.String(), amount)
	assert.Error(t, err)
	assert.EqualError(t, err, "insufficient funds")
}

func TestService_BuyItem_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)
	mockInvQuery := repository.NewMockInvQuery(ctrl)

	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery)
	mockDAO.EXPECT().NewInvQuery(gomock.Any()).Return(mockInvQuery)

	service := transaction.New(mockDAO)

	ctx := &gin.Context{}
	userID := uuid.New()
	itemType := "t-shirt"
	quantity := 1

	ctx.Set("userId", userID)

	mockInfoQuery.EXPECT().GetUserInfo(userID).Return(entity.Info{Coins: 200}, nil)
	mockInfoQuery.EXPECT().UpdateCoins(userID, 120).Return(nil)
	mockInvQuery.EXPECT().InsertInventory(userID, itemType, quantity).Return(nil)

	err := service.BuyItem(ctx, itemType, quantity)
	assert.NoError(t, err)
}

func TestService_BuyItem_InsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)

	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery)

	service := transaction.New(mockDAO)

	ctx := &gin.Context{}
	userID := uuid.New()
	itemType := "t-shirt"
	quantity := 1

	ctx.Set("userId", userID)

	mockInfoQuery.EXPECT().GetUserInfo(userID).Return(entity.Info{Coins: 50}, nil)

	err := service.BuyItem(ctx, itemType, quantity)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
}
