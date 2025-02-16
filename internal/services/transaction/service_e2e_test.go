//go:build integration
// +build integration

package transaction_test

import (
	"auth/internal/api/transaction"
	"auth/internal/repository"
	"auth/internal/repository/entity"
	"auth/internal/services/transaction"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestE2E_BuyItem(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)
	mockInvQuery := repository.NewMockInvQuery(ctrl)

	service := transaction.New(mockDAO)
	handler := api.New(service)

	userID := uuid.New()

	router.POST("/api/auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"token": "test-token"})
	})
	router.GET("/api/buy/:item", func(c *gin.Context) {
		c.Set("userId", userID)
		handler.BuyItem(c)
	})

	authResp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer([]byte(`{"username": "testuser", "password": "testpass"}`)))
	//req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(authResp, req)
	assert.Equal(t, http.StatusOK, authResp.Code)

	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery).Times(1)
	mockDAO.EXPECT().NewInvQuery(gomock.Any()).Return(mockInvQuery).Times(1)

	mockInfoQuery.EXPECT().GetUserInfo(userID).Return(entity.Info{Coins: 1000}, nil).Times(1)
	mockInfoQuery.EXPECT().UpdateCoins(userID, 1000-80).Return(nil).Times(1)
	mockInvQuery.EXPECT().InsertInventory(userID, "t-shirt", 1).Return(nil).Times(1)

	item := "t-shirt"
	buyResp := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/buy/"+item, nil)
	req.Header.Set("Authorization", "Bearer test-token")
	router.ServeHTTP(buyResp, req)

	assert.Equal(t, http.StatusOK, buyResp.Code)
	var response map[string]string
	json.Unmarshal(buyResp.Body.Bytes(), &response)
	assert.Equal(t, "Item purchased successfully", response["message"])
}

func TestE2E_TransferCoins(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockUserQuery := repository.NewMockUserQuery(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)
	mockTransQuery := repository.NewMockTransQuery(ctrl)

	service := transaction.New(mockDAO)
	handler := api.New(service)

	fromUser := uuid.New()
	toUser := uuid.New()

	router.POST("/api/auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"token": "test-token"})
	})
	router.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("userId", fromUser)
		handler.TransferCoins(c)
	})

	authResp1 := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer([]byte(`{"username": "user1", "password": "pass1"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(authResp1, req)
	assert.Equal(t, http.StatusOK, authResp1.Code)

	authResp2 := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth", bytes.NewBuffer([]byte(`{"username": "user2", "password": "pass2"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(authResp2, req)
	assert.Equal(t, http.StatusOK, authResp2.Code)

	mockDAO.EXPECT().NewUserQuery(gomock.Any()).Return(mockUserQuery).Times(1)

	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery).Times(1).Do(func(ctx context.Context) {
		fmt.Println("NewInfoQuery called")
	})

	mockDAO.EXPECT().NewTransQuery(gomock.Any()).Return(mockTransQuery).Times(1)

	mockUserQuery.EXPECT().GetUserByUsername("user2").Return(entity.User{Id: toUser}, nil).Times(1)

	mockInfoQuery.EXPECT().GetUserInfo(fromUser).Return(entity.Info{Coins: 100}, nil).Times(1)
	mockInfoQuery.EXPECT().GetUserInfo(toUser).Return(entity.Info{Coins: 50}, nil).Times(1)

	mockInfoQuery.EXPECT().UpdateCoins(fromUser, 50).Return(nil).Times(1)
	mockInfoQuery.EXPECT().UpdateCoins(toUser, 100).Return(nil).Times(1)

	mockTransQuery.EXPECT().InsertTransaction(fromUser, toUser, 50).Return(nil).Times(1)

	transferResp := httptest.NewRecorder()
	transferReqBody, _ := json.Marshal(map[string]interface{}{
		"to_user": "user2",
		"amount":  50,
	})
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(transferReqBody))
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(transferResp, req)

	assert.Equal(t, http.StatusOK, transferResp.Code)
	var response map[string]string
	json.Unmarshal(transferResp.Body.Bytes(), &response)
	assert.Equal(t, "Coins transferred successfully", response["message"])
}

func TestE2E_BuyItem_InsufficientFunds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)
	mockInvQuery := repository.NewMockInvQuery(ctrl)

	service := transaction.New(mockDAO)
	handler := api.New(service)

	userID := uuid.New()

	router.POST("/api/auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"token": "test-token"})
	})
	router.GET("/api/buy/:item", func(c *gin.Context) {
		c.Set("userId", userID)
		handler.BuyItem(c)
	})

	authResp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer([]byte(`{"username": "testuser", "password": "testpass"}`)))
	router.ServeHTTP(authResp, req)
	assert.Equal(t, http.StatusOK, authResp.Code)
	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery).Times(1)
	mockDAO.EXPECT().NewInvQuery(gomock.Any()).Return(mockInvQuery).Times(0)

	mockInfoQuery.EXPECT().GetUserInfo(userID).Return(entity.Info{Coins: 10}, nil).Times(1)

	item := "pink-hoody"
	buyResp := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/buy/"+item, nil)
	req.Header.Set("Authorization", "Bearer test-token")
	router.ServeHTTP(buyResp, req)

	assert.Equal(t, http.StatusBadRequest, buyResp.Code)
	var response map[string]string
	json.Unmarshal(buyResp.Body.Bytes(), &response)
	assert.Equal(t, "insufficient funds", response["error"])
}

func TestE2E_TransferCoins_InsufficientFunds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockUserQuery := repository.NewMockUserQuery(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)
	mockTransQuery := repository.NewMockTransQuery(ctrl)

	service := transaction.New(mockDAO)
	handler := api.New(service)

	fromUser := uuid.New()
	toUser := uuid.New()

	router.POST("/api/auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"token": "test-token"})
	})
	router.POST("/api/sendCoin", func(c *gin.Context) {
		c.Set("userId", fromUser)
		handler.TransferCoins(c)
	})

	authResp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", nil)
	router.ServeHTTP(authResp, req)
	assert.Equal(t, http.StatusOK, authResp.Code)

	mockDAO.EXPECT().NewUserQuery(gomock.Any()).Return(mockUserQuery).Times(1)
	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery).Times(1)
	mockDAO.EXPECT().NewTransQuery(gomock.Any()).Return(mockTransQuery).Times(0)

	mockUserQuery.EXPECT().GetUserByUsername("user2").Return(entity.User{Id: toUser}, nil).Times(1)
	mockInfoQuery.EXPECT().GetUserInfo(fromUser).Return(entity.Info{Coins: 10}, nil).Times(1)

	transferResp := httptest.NewRecorder()
	transferReqBody, _ := json.Marshal(map[string]interface{}{
		"to_user": "user2",
		"amount":  50,
	})
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(transferReqBody))
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(transferResp, req)

	assert.Equal(t, http.StatusBadRequest, transferResp.Code)
	var response map[string]string
	json.Unmarshal(transferResp.Body.Bytes(), &response)
	assert.Equal(t, "insufficient funds", response["error"])
}
