package info_test

import (
	"auth/internal/repository"
	"auth/internal/repository/entity"
	"auth/internal/services/info"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_GetInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)
	mockTransQuery := repository.NewMockTransQuery(ctrl)
	mockInvQuery := repository.NewMockInvQuery(ctrl)

	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery)
	mockDAO.EXPECT().NewTransQuery(gomock.Any()).Return(mockTransQuery)
	mockDAO.EXPECT().NewInvQuery(gomock.Any()).Return(mockInvQuery)

	service := info.New(mockDAO)

	ctx := &gin.Context{}
	userID := uuid.New()
	ctx.Set("userId", userID)

	mockInfoQuery.EXPECT().GetUserInfo(userID).Return(entity.Info{Coins: 500}, nil)
	mockTransQuery.EXPECT().GetSentTransactions(userID).Return([]entity.Transaction{}, nil)
	mockTransQuery.EXPECT().GetReceivedTransactions(userID).Return([]entity.Transaction{}, nil)
	mockInvQuery.EXPECT().GetUserInventory(userID).Return([]entity.Inventory{}, nil)

	infoData, err := service.GetInfo(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 500, infoData.Coins)
}

func TestService_GetInfo_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)

	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery)

	service := info.New(mockDAO)

	ctx := &gin.Context{}
	userID := uuid.New()
	ctx.Set("userId", userID)

	mockInfoQuery.EXPECT().GetUserInfo(userID).Return(entity.Info{}, errors.New("user not found"))

	infoData, err := service.GetInfo(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Equal(t, 0, infoData.Coins)
}
